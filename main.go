package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type viaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type brasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type resultado struct {
	endereco string
	erro     error
}

func main() {
	for _, cep := range os.Args[1:] {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		ch := make(chan resultado, 2)

		go buscaViaCEP(ctx, cep, ch)
		go buscaBrasilAPI(ctx, cep, ch)

		for i := 0; i < 2; i++ {
			select {
			case res := <-ch:
				if res.erro == nil {
					fmt.Printf(res.endereco)
					return
				} else {
					fmt.Println(res.erro)
				}
			case <-ctx.Done():
				fmt.Println("Erro: timeout - tente novamente mais tarde")
				return
			}
		}

	}
}

func buscaViaCEP(ctx context.Context, cep string, ch chan<- resultado) {
	//http://viacep.com.br/ws/{cep}/json/
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := novaRequisicao(ctx, url)
	if err != nil {
		ch <- resultado{"", err}
		return
	}

	var dados viaCEP
	err = json.Unmarshal(req, &dados)
	if err != nil {
		ch <- resultado{"", fmt.Errorf("Erro ao fazer parse do retorno: %v\n", err)}
	}

	endereco := fmt.Sprintf("Endereço: %s, %s - %s/%s - CEP: %s Fonte: Via CEP", dados.Logradouro, dados.Bairro, dados.Localidade, dados.Uf, dados.Cep)

	ch <- resultado{endereco, nil}
}

func buscaBrasilAPI(ctx context.Context, cep string, ch chan<- resultado) {
	//https://brasilapi.com.br/api/cep/v1/{cep}
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := novaRequisicao(ctx, url)
	if err != nil {
		ch <- resultado{"", err}
		return
	}

	var dados brasilAPI
	err = json.Unmarshal(req, &dados)
	if err != nil {
		ch <- resultado{"", fmt.Errorf("Erro ao fazer parse do retorno: %v\n", err)}
	}

	endereco := fmt.Sprintf("Endereço: %s, %s - %s/%s - CEP: %s Fonte: Brasil API", dados.Street, dados.Neighborhood, dados.City, dados.State, dados.Cep)

	ch <- resultado{endereco, nil}
}

func novaRequisicao(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Erro ao preparar a requisição: %v\n", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Erro ao realizar a requisição: %v\n", err)

	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Erro ao realizar a requisição, status HTTP inválido: %d", resp.StatusCode)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler o retorno: %v\n", err)
	}

	return res, nil
}
