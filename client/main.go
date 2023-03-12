package main

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	url := "http://localhost:8080"

	// Criando um WaitGroup para esperar todas as goroutines terminarem
	var wg sync.WaitGroup

	// Número de goroutines que vamos executar em paralelo
	numGoroutines := 30

	// Adicionando o número de goroutines ao WaitGroup
	wg.Add(numGoroutines)

	// Executando as goroutines
	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Quando a goroutine terminar, sinalize ao WaitGroup
			defer wg.Done()

			// Criando um novo cliente HTTP
			client := &http.Client{}

			// Criando uma nova requisição GET
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println("Erro ao criar a requisição:", err)
				return
			}

			// Enviando a requisição ao servidor
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Erro ao enviar a requisição:", err)
				return
			}

			// Verificando o status code da resposta
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Código de status inválido: %d\n", resp.StatusCode)
				return
			}

			fmt.Println(resp.StatusCode)

			// Lendo o corpo da resposta
			defer resp.Body.Close()
			// ... faça algo com o corpo da resposta aqui ...
		}()
	}

	// Esperando todas as goroutines terminarem
	wg.Wait()

	// Todas as goroutines terminaram
	fmt.Println("Todas as goroutines terminaram!")
}
