package tax

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/ztrixack/assessment-tax/internal/modules/database"
	"github.com/ztrixack/assessment-tax/internal/modules/logger"
)

func setup(t *testing.T) (*service, sqlmock.Sqlmock, func()) {
	log := logger.NewMockLogger()
	db, mock, err := database.NewMockDB()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	svr := New(log, db)

	return svr, mock, func() {
		db.Close()
	}
}

func TestCalculateStepTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		lower       float64
		upper       float64
		rate        float64
		expectedTax float64
	}{
		{"Zero rate", 100000, 0, 150000, 0, 0},
		{"Within first bracket", 200000, 150000, 500000, 0.10, 5000},
		{"Within second bracket", 750000, 500000, 1000000, 0.15, 37500},
		{"Within third bracket", 1500000, 1000000, 2000000, 0.20, 100000},
		{"Within fourth bracket", 3000000, 2000000, 3000000, 0.35, 350000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax := calculateStepTax(tt.income, tt.lower, tt.upper, tt.rate)
			assert.Equal(t, tt.expectedTax, gotTax)
		})
	}
}

func TestCalculateProgressiveTax(t *testing.T) {
	tests := []struct {
		name        string
		income      float64
		expectedTax float64
		expectedErr error
	}{
		{
			name:        "Story: EXP01",
			income:      440000,
			expectedTax: 29000,
			expectedErr: nil,
		},
		{
			name:        "Story: EXP03",
			income:      340000,
			expectedTax: 19000,
			expectedErr: nil,
		},
		{
			name:        "Negative income",
			income:      -100,
			expectedTax: 0,
			expectedErr: ErrNegativeIncome,
		},
		{
			name:        "Zero income",
			income:      0,
			expectedTax: 0,
			expectedErr: nil,
		},
		{
			name:        "Income within first bracket (0% tax)",
			income:      150000,
			expectedTax: 0,
			expectedErr: nil,
		},
		{
			name:        "Income within second bracket on the lower end (10% tax)",
			income:      150001,
			expectedTax: 0.1,
			expectedErr: nil,
		},
		{
			name:        "Income within second bracket on the upper end (10% tax)",
			income:      500000,
			expectedTax: 35000,
			expectedErr: nil,
		},
		{
			name:        "Income within third bracket on the lower end (15% tax)",
			income:      500001,
			expectedTax: 35000.15,
			expectedErr: nil,
		},
		{
			name:        "Income within third bracket on the upper end (15% tax)",
			income:      1000000,
			expectedTax: 110000,
			expectedErr: nil,
		},
		{
			name:        "Income within fourth bracket on the lower end (20% tax)",
			income:      1000001,
			expectedTax: 110000.2,
			expectedErr: nil,
		},
		{
			name:        "Income within fourth bracket on the upper end (20% tax)",
			income:      2000000,
			expectedTax: 310000,
			expectedErr: nil,
		},
		{
			name:        "Income within fifth bracket on the lower end (35% tax)",
			income:      2000001,
			expectedTax: 310000.35,
			expectedErr: nil,
		},
		{
			name:        "Income within fifth bracket (35% tax)",
			income:      5000000,
			expectedTax: 1360000,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTax, gotErr := calculateProgressiveTax(tt.income)
			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedTax, gotTax)
		})
	}
}

func TestCalculateAllowances(t *testing.T) {
	defaultMockBehavior := func(mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"personal", "donation"}).AddRow(60000, 100000)
		mock.ExpectPrepare("SELECT personal, donation FROM allowances").ExpectQuery().WillReturnRows(rows)
	}

	tests := []struct {
		name           string
		mockBehavior   func(sqlmock.Sqlmock)
		allowances     []Allowance
		expectedResult float64
		wantErr        bool
	}{
		{
			name:           "Story: EXP01",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{{Type: "donation", Amount: 0}},
			expectedResult: 60000,
			wantErr:        false,
		},
		{
			name:           "Story: EXP03",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{{Type: Donation, Amount: 200000}},
			expectedResult: 60000 + 100000,
			wantErr:        false,
		},
		{
			name:           "No allowances",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{},
			expectedResult: 60000,
			wantErr:        false,
		},
		{
			name:           "All minimum values",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{{Type: Donation, Amount: 0}},
			expectedResult: 60000,
			wantErr:        false,
		},
		{
			name:           "All maximum values",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{{Type: Donation, Amount: 100000}},
			expectedResult: 60000 + 100000,
			wantErr:        false,
		},
		{
			name:         "Negative amounts",
			mockBehavior: defaultMockBehavior,
			allowances:   []Allowance{{Type: Donation, Amount: -50}},
			wantErr:      true,
		},
		{
			name:           "Above maximum limits",
			mockBehavior:   defaultMockBehavior,
			allowances:     []Allowance{{Type: Donation, Amount: 100001}},
			expectedResult: 60000 + 100000,
			wantErr:        false,
		},
		{
			name:         "Multi allowances and below maximum limits",
			mockBehavior: defaultMockBehavior,
			allowances: []Allowance{
				{Type: Donation, Amount: 30000},
				{Type: Donation, Amount: 30000},
			},
			expectedResult: 60000 + 60000,
			wantErr:        false,
		},
		{
			name:         "Multi allowances and above maximum limits",
			mockBehavior: defaultMockBehavior,
			allowances: []Allowance{
				{Type: Donation, Amount: 60000},
				{Type: Donation, Amount: 80000},
			},
			expectedResult: 60000 + 100000,
			wantErr:        false,
		},
		{
			name:         "Unknown allowance type",
			mockBehavior: defaultMockBehavior,
			allowances:   []Allowance{{Type: "uknnown", Amount: 30000}},
			wantErr:      true,
		},
		{
			name: "Error with database",
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("SELECT personal, donation FROM allowances").ExpectQuery().WillReturnError(errors.New("some error"))
			},
			allowances: []Allowance{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr, mock, close := setup(t)
			defer close()

			tt.mockBehavior(mock)

			result, err := svr.calculateAllowances(tt.allowances)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
