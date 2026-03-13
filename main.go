package main

import (
	"CoffeeShopOnGo/service"
	"CoffeeShopOnGo/api"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	ginSwagger "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "CoffeeShopOnGo/docs"
)

type App struct {
	svc    *service.Service
	report *service.ReportService
	reader *bufio.Reader
}

func NewApp(svc *service.Service, report *service.ReportService) *App {
	return &App{svc: svc,
		report: report,
		reader: bufio.NewReader(os.Stdin),
	}
}

func (app *App) readString(prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := app.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func (app *App) readInt(prompt string) (int, error) {
	str, err := app.readString(prompt)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func (app *App) readFloat(prompt string) (float64, error) {
	str, err := app.readString(prompt)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(str, 64)
}

func (app *App) readBool(prompt string) (bool, error) {
	str, err := app.readString(prompt + " (y/n): ")
	if err != nil {
		return false, err
	}
	return strings.ToLower(str) == "y", nil
}

func ClearConsole() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Меню товаров

func (app *App) productMenu() {
	for {
		ClearConsole()
		fmt.Println("\n=== ТОВАРЫ ===")
		fmt.Println("1. Список товаров")
		fmt.Println("2. Добавить товар")
		fmt.Println("3. Изменить цену")
		fmt.Println("0. Назад")

		choice, err := app.readString(">>> ")
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case "1":
			app.showProducts()
		case "2":
			app.addProduct()
		case "3":
			app.changePrice()
		case "0":
			return
		default:
			fmt.Println("Неверный выбор")
		}
		app.pause()
	}
}

func (app *App) showProducts() {
	products, err := app.svc.ShowAllProducts()
	if err != nil {
		fmt.Printf("Ошибка загрузки: %v\n", err)
		return
	}

	if len(products) == 0 {
		fmt.Println("Товаров нет")
		return
	}
	fmt.Println("\nСПИСОК ТОВАРОВ:")
	for _, p := range products {
		fmt.Printf("ID: %d | %s | Цена: %.2f | Категория: %s | НДС: %d%%\n", p.ID, p.Name, p.Price, p.Category, p.VAT)
	}
}

func (app *App) addProduct() {
	fmt.Println("\nНОВЫЙ ТОВАР")

	name, err := app.readString("Название: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	price, err := app.readFloat("Цена: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	category, err := app.readString("Категория: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	vat, err := app.readInt("НДС (%): ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	isActive, err := app.readBool("Активен")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}

	product, err := app.svc.AddNewProduct(name, price, category, vat, isActive)
	if err != nil {
		fmt.Printf("Ошибка создания товара: %v\n", err)
		return
	}

	fmt.Printf("Товар '%s' добавлен с ID %d\n", product.Name, product.ID)
}

func (app *App) changePrice() {
	fmt.Println("\nИЗМЕНЕНИЕ ЦЕНЫ")

	id, err := app.readInt("ID товара: ")
	if err != nil {
		fmt.Println("Неверный ID")
		return
	}

	// Сначала покажем текущий товар
	product, err := app.svc.GetProduct(id)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if product == nil {
		fmt.Println("Товар не найден")
		return
	}

	fmt.Printf("Текущий товар: %s, цена: %.2f\n", product.Name, product.Price)

	newPrice, err := app.readFloat("Новая цена: ")
	if err != nil {
		fmt.Println("Неверный формат цены")
		return
	}

	updated, err := app.svc.ChangePrice(id, newPrice)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	fmt.Printf("Цена товара '%s' изменена на %.2f\n", updated.Name, updated.Price)
}

// Персонал
func (app *App) staffMenu() {
	for {
		ClearConsole()
		fmt.Println("\n=== ПЕРСОНАЛ ===")
		fmt.Println("1. Список")
		fmt.Println("2. Нанять")
		fmt.Println("3. Уволить")
		fmt.Println("0. Назад")

		choice, err := app.readString(">>> ")
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case "1":
			app.showStaff()
		case "2":
			app.addStaff()
		case "3":
			app.fireStaff()
		case "0":
			return
		default:
			fmt.Println("Неверный выбор")
		}
		app.pause()
	}
}

func (app *App) showStaff() {
	staff, err := app.svc.ShowAllStaff()
	if err != nil {
		fmt.Printf("Ошибка загрузки: %v\n", err)
		return
	}
	if len(staff) == 0 {
		fmt.Println("Работников нет")
		return
	}
	fmt.Println("\nСПИСОК РАБОТНИКОВ:")
	for _, s := range staff {
		fmt.Printf("ID: %d | ФИО: %s | ЗП: %.2f | Должность: %s | Возраст: %d%%\n", s.ID, s.FullName, s.Salary, s.Position, s.Age)
	}
}

func (app *App) addStaff() {
	fmt.Println("\nНОВЫЙ РАБОТНИК")

	name, err := app.readString("ФИО: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	sal, err := app.readFloat("Зарплата: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	pos, err := app.readString("Должность: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	age, err := app.readInt("Возраст: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}

	staff, err := app.svc.AddNewStaff(name, sal, age, pos)
	if err != nil {
		fmt.Printf("Ошибка создания сотрудника: %v\n", err)
		return
	}

	fmt.Printf("Сотрудник '%s' добавлен с ID %d\n", staff.FullName, staff.ID)
}

func (app *App) fireStaff() {
	id, err := app.readInt("ID работника: ")
	if err != nil {
		fmt.Println("Ошибка ввода")
		return
	}
	er := app.svc.FireStaff(id)
	if er != nil {
		fmt.Printf("Ошибка увольнения сотрудника: %v\n", er)
		return
	}
	fmt.Println("Сотрудник уволен")

}

// Продажи
func (app *App) salesMenu() {
	for {
		ClearConsole()
		fmt.Println("\n=== Продажи ===")
		fmt.Println("1. Список")
		fmt.Println("2. Добавить")
		fmt.Println("3. Общая выручка")
		fmt.Println("0. Назад")

		choice, err := app.readString(">>> ")
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case "1":
			app.showAllSales()
		case "2":
			app.addSale()
		case "3":
			app.showRevenue()
		case "0":
			return
		default:
			fmt.Println("Неверный выбор")
		}
		app.pause()
	}
}

func (app *App) showAllSales() {
    sales, err := app.svc.GetAllSales()
    if err != nil {
        fmt.Printf("Ошибка загрузки: %v\n", err)
        return
    }
    
    if len(sales) == 0 {
        fmt.Println("Продаж нет")
        return
    }
    
    fmt.Println("\nСПИСОК ПРОДАЖ:")
    for _, s := range sales {
        fmt.Printf("ID: %d | Сумма: %.2f | Кол-во: %d | Дата: %s\n",
            s.ID, s.Amount, s.Quantity, s.CreatedAt.Format("02.01.2006 15:04"))
    }
}

func (app *App) addSale() {
    fmt.Println("\nНОВАЯ ПРОДАЖА")
    
    amount, err := app.readFloat("Сумма: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    productId, err := app.readInt("ID товара: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    staffId, err := app.readInt("ID сотрудника: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    quantity, err := app.readInt("Количество: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    sale, err := app.svc.AddNewSale(amount, productId, staffId, quantity)
    if err != nil {
        fmt.Printf("Ошибка создания продажи: %v\n", err)
        return
    }
    
    fmt.Printf("Продажа №%d добавлена\n", sale.ID)
}

func (app *App) showRevenue() {
	rev, err := app.svc.Revenue()
	if err != nil {
		fmt.Printf("Ошибка подсчета прибыли: %v\n", err)
		return
	}
	fmt.Printf("Общая выручка: %.2f рублей", rev)
}

// Склад
func (app *App) warehouseMenu() {
    for {
        ClearConsole()
        fmt.Println("\n=== СКЛАД ===")
        fmt.Println("1. Список позиций")
        fmt.Println("2. Добавить позицию")
        fmt.Println("3. Изменить остаток")
        fmt.Println("4. Проверить остатки (менее 10)")
        fmt.Println("0. Назад")
        
        choice, err := app.readString(">>> ")
        if err != nil {
            fmt.Println("Ошибка ввода:", err)
            continue
        }
        
        switch choice {
        case "1":
            app.showWarehouse()
        case "2":
            app.addWarehouseItem()
        case "3":
            app.changeStock()
        case "4":
            app.checkLowStock()
        case "0":
            return
        default:
            fmt.Println("Неверный выбор")
        }
        app.pause()
    }
}

func (app *App) showWarehouse() {
    items, err := app.svc.GetAllWHItems()
    if err != nil {
        fmt.Printf("Ошибка загрузки: %v\n", err)
        return
    }
    
    if len(items) == 0 {
        fmt.Println("Склад пуст")
        return
    }
    
    fmt.Println("\nСОДЕРЖИМОЕ СКЛАДА:")
    for _, w := range items {
        // Подгружаем продукт для имени
        var product *service.Product
        product, err = app.svc.GetProduct(w.ProductId)
        productName := "Неизвестно"
        if err == nil && product != nil {
            productName = product.Name
        }
        
        fmt.Printf("ID: %d | Товар: %s | Остаток: %d | Статус: %s\n",
            w.ID, productName, w.Stock, w.Status)
    }
}

func (app *App) addWarehouseItem() {
    fmt.Println("\nНОВАЯ ПОЗИЦИЯ НА СКЛАДЕ")
    
    productId, err := app.readInt("ID товара: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    // Проверим, что товар существует
    product, err := app.svc.GetProduct(productId)
    if err != nil || product == nil {
        fmt.Println("Товар не найден")
        return
    }
    
    stock, err := app.readInt("Количество: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    status, err := app.readString("Статус (например, 'доступен', 'закончился'): ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    item, err := app.svc.AddNewWHItem(productId, stock, status)
    if err != nil {
        fmt.Printf("Ошибка добавления: %v\n", err)
        return
    }
    
    fmt.Printf("Позиция добавлена с ID %d\n", item.ID)
}

func (app *App) changeStock() {
    fmt.Println("\nИЗМЕНЕНИЕ ОСТАТКА")
    
    id, err := app.readInt("ID позиции на складе: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    fmt.Println("Выберите действие:")
    fmt.Println("1. Увеличить")
    fmt.Println("2. Уменьшить")
    
    action, err := app.readString(">>> ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    amount, err := app.readInt("Количество: ")
    if err != nil {
        fmt.Println("Ошибка ввода")
        return
    }
    
    var method string
    switch action {
    case "1":
        method = "increase"
    case "2":
        method = "decrease"
    default:
        fmt.Println("Неверное действие")
        return
    }
    
    err = app.svc.ChangeStock(id, amount, method)
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    
    fmt.Println("Остаток обновлен")
}

func (app *App) checkLowStock() {
    alert, err := app.report.LowStockAlert(10)  // порог 10
    if err != nil {
        fmt.Printf("Ошибка: %v\n", err)
        return
    }
    fmt.Println(alert)
}

// Отчеты
func (app *App) reportMenu() {
	for {
		ClearConsole()
		fmt.Println("\n=== ОТЧЕТ ===")
		fmt.Println("1. Получить отчет")
		fmt.Println("2. Скачать отчет")
		fmt.Println("0. Назад")

		choice, err := app.readString(">>> ")
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case "1":
			app.showReport()
		case "2":
			app.downloadReport()
		case "0":
			return
		default:
			fmt.Println("Неверный выбор")
		}
		app.pause()
	}
}

func (app *App) showReport() {
	rep, err := app.report.Report()
	if err != nil {
		fmt.Printf("Ошибка при создании отчета: %v\n", err)
		return
	}
	fmt.Println(rep)
}

func (app *App) downloadReport() {
	rep, err := app.report.Report()
	if err != nil {
		fmt.Printf("Ошибка при создании отчета: %v\n", err)
		return
	}
	currentDateTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("report_%s.txt", currentDateTime)

	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(rep); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Отчет скачан в %s", fileName)

}

// Вспомогательное

func (app *App) pause() {
	fmt.Println("\nНажмите Enter чтобы продолжить...")
	app.reader.ReadString('\n')
}

// Главное меню

func (app *App) run() {
	for {
		ClearConsole()
		fmt.Println("\n=== КОФЕЙНЯ ===")
		fmt.Println("1. Товары")
		fmt.Println("2. Персонал")
		fmt.Println("3. Продажи")
		fmt.Println("4. Склад")
		fmt.Println("5. Отчеты")
		fmt.Println("0. Выход")

		choice, err := app.readString(">>> ")
		if err != nil {
			fmt.Println("Ошибка ввода, попробуйте снова")
			continue
		}

		switch choice {
		case "1":
			app.productMenu()
		case "2":
			app.staffMenu()
		case "3":
			app.salesMenu()
		case "4":
			app.warehouseMenu()
		case "5":
			app.reportMenu()
		case "0":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный выбор. Введите 0-5")
		}
	}
}

// @title           Coffee Shop API
// @version         1.0
// @description     Сервис для управления товарами в кофейне.
// @host      localhost:8080
// @BasePath  /
func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	repo := service.NewRepository(db)

	if err := repo.AutoMigrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)

	}

	svc := service.NewService(repo)
	reportSvc := service.NewReportService(svc)

	app := NewApp(svc, reportSvc)
	app.run()

	router := api.SetupRouter(svc)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    log.Println("Сервер запущен на :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
}
