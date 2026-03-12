package service

import (
	"errors"
	"fmt"
	"sort"
	"time" 
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}


// Сотрудники
func (s *Service) GetBestWorker() (*Staff, error) {
	var staffList []Staff
	err := s.repo.db.Preload("Sales").Find(&staffList).Error
	if err != nil {
		return nil, err
	}
	var best *Staff
	maxSales := 0
	for i := range staffList {
		if len(staffList[i].Sales) > maxSales {
			maxSales = len(staffList[i].Sales)
			best = &staffList[i]
		}
	}
	if best == nil {
		return nil, errors.New("Нет сотрудников")
	}
	return best, nil
}

func (s *Service) ShowAllStaff() ([]Staff, error) {
	staff, err := s.repo.Staff.GetAll()
	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *Service) GetOneStaff(id int) (*Staff, error) {
	staff, err := s.repo.Staff.Get(id)
	if err != nil {
		return nil, err
	}
	return staff, nil
}


func (s *Service) AddNewStaff(fullName string, salary float64, age int, position string) (*Staff, error) {
	staffEntity := &Staff{
		FullName: fullName,
		Salary:   salary,
		Age:      age,
		Position: position,
	}
	staff, err := s.repo.Staff.Create(staffEntity)
	if err != nil {
		return nil, fmt.Errorf("добавление сотрудника: %w", err)
	}
	return staff, nil
}

func (s *Service) FireStaff(id int) error {
	success, err := s.repo.Staff.Delete(id)
	if err != nil {
		return fmt.Errorf("удаление сотрудника №%d: %w", id, err)
	}
	if !success {
		return fmt.Errorf("сотрудник %d не найден", id)
	}
	return nil
}

// Товары
func (s *Service) GetMostSoldProduct() (*Product, error) {
	var productList []Product
	err := s.repo.db.Preload("Sales").Find(&productList).Error
	if err != nil {
		return nil, err
	}
	var most *Product
	maxSales := 0
	for i := range productList {
		if len(productList[i].Sales) > maxSales {
			maxSales = len(productList[i].Sales)
			most = &productList[i]
		}
	}
	if most == nil {
		return nil, errors.New("Нет товаров")
	}
	return most, nil
}

func (s *Service) ShowAllProducts() ([]Product, error) {
	products, err := s.repo.Product.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *Service) GetProduct(id int) (*Product, error) {
	p, err := s.repo.Product.Get(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}


func (s *Service) AddNewProduct(name string, price float64, category string, vat int, isActive bool) (*Product, error) {
	product := &Product{
		Name:     name,
		Price:    price,
		Category: category,
		VAT:      vat,
		IsActive: isActive,
	}
	p, err := s.repo.Product.Create(product)
	if err != nil {
		return nil, fmt.Errorf("добавление товара: %w", err)
	}
	return p, nil
}

func (s *Service) ChangePrice(id int, amount float64) (*Product, error) {
	product, err := s.repo.Product.Get(id)
	if err != nil {
		return nil, fmt.Errorf("получение товара: %w", err)
	}
	if product == nil {
		return nil, errors.New("товар не найден")
	}
	product.Price = amount
	p, err := s.repo.Product.Update(id, product)
	if err != nil {
		return nil, fmt.Errorf("обновление цены: %w", err)
	}
	return p, nil
}

// Продажи
func (s *Service) GetMostCostSale() (*Sale, error) {
	sales, err := s.repo.Sale.GetAll()
	if err != nil {
		return nil, err
	}
	if len(sales) == 0 {
		return nil, nil
	}
	sort.Slice(sales, func(i, j int) bool {
		return sales[i].Amount > sales[j].Amount
	})
	return &sales[0], nil
}


func (s *Service) GetAllSales() ([]Sale, error) {
	sales, err := s.repo.Sale.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение продаж: %w", err)
	}
	return sales, nil
}


func (s *Service) AddNewSale(amount float64, productId, staffId, quantity int) (*Sale, error) {
	// Проверяем существование сотрудника
	staff, err := s.repo.Staff.Get(staffId)
	if err != nil || staff == nil {
		return nil, errors.New("сотрудник не найден")
	}
	
	// Проверяем существование товара
	product, err := s.repo.Product.Get(productId)
	if err != nil || product == nil {
		return nil, errors.New("товар не найден")
	}

	sale := &Sale{
		Amount:    amount,
		CreatedAt: time.Now(),
		ProductId: productId,
		StaffId:   staffId,
		Quantity:  quantity,
	}
	
	result, err := s.repo.Sale.Create(sale)
	if err != nil {
		return nil, fmt.Errorf("создание продажи: %w", err)
	}
	return result, nil
}

// Revenue - общая выручка
func (s *Service) Revenue() (float64, error) {
	sales, err := s.repo.Sale.GetAll()
	if err != nil {
		return 0, fmt.Errorf("получение выручки: %w", err)
	}
	var total float64
	for _, sale := range sales {
		total += sale.Amount
	}
	return total, nil
}

// Склад
func (s *Service) GetAllWHItems() ([]Warehouse, error) {
	items, err := s.repo.Warehouse.GetAll()
	if err != nil {
		return nil, fmt.Errorf("получение склада: %w", err)
	}
	return items, nil
}

func (s *Service) AddNewWHItem(productId int, stock int, status string) (*Warehouse, error) {
	// Проверяем, что продукт существует
	product, err := s.repo.Product.Get(productId)
	if err != nil || product == nil {
		return nil, errors.New("продукт не найден")
	}

	wh := &Warehouse{
		ProductId: productId,
		Stock:     stock,
		Status:    status,
	}
	
	result, err := s.repo.Warehouse.Create(wh)
	if err != nil {
		return nil, fmt.Errorf("добавление на склад: %w", err)
	}
	return result, nil
}

func (s *Service) GetMostInStock() (*Warehouse, error) {
	items, err := s.repo.Warehouse.GetAll()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	
	sort.Slice(items, func(i, j int) bool {
		return items[i].Stock > items[j].Stock
	})
	
	var result Warehouse
	err = s.repo.db.Preload("Product").First(&result, items[0].ID).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) ChangeStock(id int, amount int, method string) error {
	wh, err := s.repo.Warehouse.Get(id)
	if err != nil {
		return fmt.Errorf("получение позиции: %w", err)
	}
	if wh == nil {
		return errors.New("позиция не найдена")
	}

	switch method {
	case "increase":
		wh.Stock += amount
	case "decrease":
		if wh.Stock < amount {
			return errors.New("недостаточно товара на складе")
		}
		wh.Stock -= amount
	default:
		return errors.New("недопустимый метод")
	}

	_, err = s.repo.Warehouse.Update(id, wh)
	if err != nil {
		return fmt.Errorf("обновление склада: %w", err)
	}
	return nil
}



// ReportService - структура для отчетов

type ReportService struct {
	svc *Service // зависимость от основного сервиса
}

func NewReportService(svc *Service) *ReportService {
	return &ReportService{svc: svc}
}


func (rs *ReportService) Report() (string, error) {
	var sb strings.Builder 
	
	sb.WriteString("========== ОТЧЕТ ==========\n")
	sb.WriteString(fmt.Sprintf("Дата: %s\n\n", time.Now().Format("02.01.2006"))) 
	sb.WriteString("Краткая статистика:\n")

	// Выручка
	revenue, err := rs.svc.Revenue()
	if err != nil {
		return "", fmt.Errorf("ошибка получения выручки: %w", err)
	}
	sb.WriteString(fmt.Sprintf("Общая выручка: %.2f рублей\n", revenue))

	// Лучший работник
	bestWorker, err := rs.svc.GetBestWorker()
	if err != nil {
		sb.WriteString("Лучший работник: ошибка получения\n")
	} else if bestWorker != nil {
		// Preload уже подгрузил Sales
		sb.WriteString(fmt.Sprintf("Лучший работник: %s (%d продаж)\n", 
			bestWorker.FullName, len(bestWorker.Sales)))
	} else {
		sb.WriteString("Лучший работник: нет данных\n")
	}

	// Самая дорогая продажа
	mostExpensive, err := rs.svc.GetMostCostSale()
	if err != nil {
		sb.WriteString("Самая дорогая продажа: ошибка получения\n")
	} else if mostExpensive != nil {
		sb.WriteString(fmt.Sprintf("Самая дорогая продажа: %.2f рублей\n", mostExpensive.Amount))
	} else {
		sb.WriteString("Самая дорогая продажа: нет данных\n")
	}

	// Самый популярный товар
	mostSold, err := rs.svc.GetMostSoldProduct()
	if err != nil {
		sb.WriteString("Самый популярный товар: ошибка получения\n")
	} else if mostSold != nil {
		sb.WriteString(fmt.Sprintf("Самый популярный товар: %s\n", mostSold.Name))
	} else {
		sb.WriteString("Самый популярный товар: нет данных\n")
	}

	// Больше всего на складе
	mostStock, err := rs.svc.GetMostInStock()
	if err != nil {
		sb.WriteString("Больше всего на складе: ошибка получения\n")
	} else if mostStock != nil && mostStock.Product != nil {
		sb.WriteString(fmt.Sprintf("Больше всего на складе: %s\n", mostStock.Product.Name))
	} else {
		sb.WriteString("Склад: нет данных\n")
	}

	// Подробная информация
	sb.WriteString("\nПодробная информация:\n")
	
	// Персонал
	sb.WriteString("Персонал:\n")
	staffList, err := rs.svc.ShowAllStaff()
	if err != nil {
		sb.WriteString("  Ошибка загрузки\n")
	} else {
		for _, s := range staffList {
			// Нужно подгрузить продажи для каждого
			var staffWithSales Staff
			err := rs.svc.repo.db.Preload("Sales").First(&staffWithSales, s.ID).Error
			salesCount := 0
			if err == nil {
				salesCount = len(staffWithSales.Sales)
			}
			sb.WriteString(fmt.Sprintf("  ФИО: %s | ЗП: %.2f | Должность: %s | Возраст: %d | Продаж: %d\n",
				s.FullName, s.Salary, s.Position, s.Age, salesCount))
		}
	}

	// Товары
	sb.WriteString("\nТовары:\n")
	products, err := rs.svc.ShowAllProducts()
	if err != nil {
		sb.WriteString("  Ошибка загрузки\n")
	} else {
		for _, p := range products {
			sb.WriteString(fmt.Sprintf("  Название: %s | Цена: %.2f рублей | Категория: %s | НДС: %d%%\n",
				p.Name, p.Price, p.Category, p.VAT))
		}
	}

	// Продажи
	sb.WriteString("\nПродажи:\n")
	sales, err := rs.svc.GetAllSales()
	if err != nil {
		sb.WriteString("  Ошибка загрузки\n")
	} else {
		for _, sale := range sales {
			sb.WriteString(fmt.Sprintf("  ID: %d | Сумма: %.2f | Кол-во: %d | Дата: %s\n",
				sale.ID, sale.Amount, sale.Quantity, sale.CreatedAt.Format("02.01.2006 15:04")))
		}
	}

	// Склад
	sb.WriteString("\nСклад:\n")
	whItems, err := rs.svc.GetAllWHItems()
	if err != nil {
		sb.WriteString("  Ошибка загрузки\n")
	} else {
		for _, w := range whItems {
			// Подгружаем продукт для имени
			var whWithProduct Warehouse
			err := rs.svc.repo.db.Preload("Product").First(&whWithProduct, w.ID).Error
			productName := "Неизвестно"
			if err == nil && whWithProduct.Product != nil {
				productName = whWithProduct.Product.Name
			}
			sb.WriteString(fmt.Sprintf("  ID: %d | Товар: %s | Остаток: %d | Статус: %s\n",
				w.ID, productName, w.Stock, w.Status))
		}
	}

	return sb.String(), nil
}

// DailyReport - отчет за день
func (rs *ReportService) DailyReport(date time.Time) (string, error) {
	sales, err := rs.svc.GetAllSales()
	if err != nil {
		return "", fmt.Errorf("ошибка получения продаж: %w", err)
	}

	var total float64
	for _, sale := range sales {
		// Сравниваем даты (без времени)
		if sale.CreatedAt.Year() == date.Year() &&
			sale.CreatedAt.Month() == date.Month() &&
			sale.CreatedAt.Day() == date.Day() {
			total += sale.Amount
		}
	}

	return fmt.Sprintf("Выручка за %s: %.2f рублей", 
		date.Format("02.01.2006"), total), nil
}

// LowStockAlert - проверка остатков
func (rs *ReportService) LowStockAlert(threshold int) (string, error) {
	whItems, err := rs.svc.GetAllWHItems()
	if err != nil {
		return "", fmt.Errorf("ошибка получения склада: %w", err)
	}

	var lowStock []string
	for _, w := range whItems {
		if w.Stock < threshold {
			// Подгружаем продукт для имени
			var whWithProduct Warehouse
			err := rs.svc.repo.db.Preload("Product").First(&whWithProduct, w.ID).Error
			if err == nil && whWithProduct.Product != nil {
				lowStock = append(lowStock, whWithProduct.Product.Name)
			} else {
				lowStock = append(lowStock, fmt.Sprintf("ID:%d", w.ProductId))
			}
		}
	}

	if len(lowStock) == 0 {
		return "Склад в норме", nil
	}
	
	return fmt.Sprintf("Нехватка: %s", strings.Join(lowStock, ", ")), nil
}