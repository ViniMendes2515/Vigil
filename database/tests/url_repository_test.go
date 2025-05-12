package database

import (
	"errors"
	"testing"
	"vigil/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Mocks seguindo o padrao AAA (Arrange, Act, Assert)

// TestAddUrl verifica a função AddUrl
func TestAddUrl(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	site := "amazon.com"
	url := "https://amazon.com/produto1"
	name := "Produto 1"
	precoInicial := 100.0
	precoLimite := 90.0

	mockRepo.EXPECT().
		AddUrl(site, url, name, precoInicial, precoLimite).
		Return(nil)

	err := mockRepo.AddUrl(site, url, name, precoInicial, precoLimite)

	assert.NoError(t, err)
}

// TestAddUrl_Error verifica o tratamento de erros em AddUrl
func TestAddUrl_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	erroEsperado := errors.New("erro ao adicionar URL")

	mockRepo.EXPECT().
		AddUrl(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(erroEsperado)

	err := mockRepo.AddUrl("amazon.com", "https://amazon.com/produto1", "Produto 1", 100.0, 90.0)

	assert.Error(t, err)
	assert.Equal(t, erroEsperado, err)
}

// TestGetSiteUrls verifica a função GetSiteUrls
func TestGetSiteUrls(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	host := "amazon.com"
	expectedUrls := []string{"https://amazon.com/produto1", "https://amazon.com/produto2"}

	mockRepo.EXPECT().
		GetSiteUrls(host).
		Return(expectedUrls, nil)

	urls, err := mockRepo.GetSiteUrls(host)

	assert.NoError(t, err)
	assert.Equal(t, expectedUrls, urls)
}

// TestGetSiteUrls_Error verifica o tratamento de erros em GetSiteUrls
func TestGetSiteUrls_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockUrlRepository(ctrl)

	expectedError := errors.New("erro ao obter URLs do site")

	mockRepo.EXPECT().
		GetSiteUrls(gomock.Any()).
		Return(nil, expectedError)

	urls, err := mockRepo.GetSiteUrls("amazonW3C")
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, urls)
}

// TestGetUrl verifica a função GetUrls
func TestGetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	expectedUrl := map[string]int{
		"https://example.com/produto1": 1,
		"https://example.com/produto2": 2,
	}

	mockRepo.EXPECT().
		GetUrls().
		Return(expectedUrl, nil)

	urls, err := mockRepo.GetUrls()

	assert.NoError(t, err)
	assert.Equal(t, expectedUrl, urls)
}

// TestGetUrl_Error verifica o tratamento de erros em GetUrls
func TestGetUrl_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	expectedError := errors.New("erro ao obter URLs")

	mockRepo.EXPECT().
		GetUrls().
		Return(nil, expectedError)

	urls, err := mockRepo.GetUrls()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, urls)
}

// TestRemoveAllUrls verifica a função RemoveAllUrls
func TestRemoveAllUrls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)

	mockRepo.EXPECT().
		RemoveAllUrls().
		Return(nil)

	err := mockRepo.RemoveAllUrls()

	assert.NoError(t, err)
}

// TestRemoveAllUrls_Erro verifica o tratamento de erros em RemoveAllUrls
func TestRemoveAllUrls_Erro(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	erroEsperado := errors.New("erro ao remover URLs")

	mockRepo.EXPECT().
		RemoveAllUrls().
		Return(erroEsperado)

	err := mockRepo.RemoveAllUrls()

	assert.Error(t, err)
	assert.Equal(t, erroEsperado, err)
}

// TestRemoveUrl verifica a função RemoveUrl
func TestRemoveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	urlParaRemover := "https://example.com/produto"

	mockRepo.EXPECT().
		RemoveUrl(urlParaRemover).
		Return(nil)

	err := mockRepo.RemoveUrl(urlParaRemover)

	assert.NoError(t, err)
}

// TestRemoveUrl_Erro verifica o tratamento de erros em RemoveUrl
func TestRemoveUrl_Erro(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	erroEsperado := errors.New("erro ao remover URL")

	mockRepo.EXPECT().
		RemoveUrl(gomock.Any()).
		Return(erroEsperado)

	err := mockRepo.RemoveUrl("url-qualquer")

	assert.Error(t, err)
	assert.Equal(t, erroEsperado, err)
}

// TestRemoveUrlById verifica a função RemoveUrlById
func TestRemoveUrlById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	idParaRemover := 42

	mockRepo.EXPECT().
		RemoveUrlById(idParaRemover).
		Return(nil)

	err := mockRepo.RemoveUrlById(idParaRemover)

	assert.NoError(t, err)
}

// TestRemoveUrlById_Erro verifica o tratamento de erros em RemoveUrlById
func TestRemoveUrlById_Erro(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUrlRepository(ctrl)
	erroEsperado := errors.New("erro ao remover URL por ID")

	mockRepo.EXPECT().
		RemoveUrlById(gomock.Any()).
		Return(erroEsperado)

	err := mockRepo.RemoveUrlById(123)

	assert.Error(t, err)
	assert.Equal(t, erroEsperado, err)
}
