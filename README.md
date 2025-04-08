# comunica_serial

`comunica_serial` é um módulo Go para comunicação serial, projetado para facilitar o envio de dados através de portas seriais. Este repositório contém um exemplo de uso e a implementação de um sistema de envio de dados periódicos.

## Instalação

Para usar este módulo em seu próprio código, você precisa ter o Go instalado em sua máquina. Certifique-se de que sua versão do Go seja compatível com o módulo (veja o arquivo `go.mod` para detalhes).

1. Adicione o módulo `comunica_serial` ao seu projeto usando o comando `go get`:

    ```bash
    go get github.com/TsukiGva2/comunica_serial
    ```

2. Importe o módulo em seu código Go:

    ```go
    import "github.com/TsukiGva2/comunica_serial"
    ```

3. Certifique-se de executar o comando abaixo para atualizar as dependências do seu projeto:

    ```bash
    go mod tidy
    ```

## Exemplo de Uso

Abaixo está um exemplo básico de como usar o módulo `comunica_serial` para enviar dados através de uma porta serial:

```go
package main

import (
    "log"
    "time"
)

func main() {
    // Inicializa o SerialSender com uma taxa de baud de 115200
    sender, err := NewSerialSender(115200)
    if err != nil {
        log.Fatalf("Falha ao inicializar o SerialSender: %v", err)
    }
    defer sender.Close()

    // Cria uma instância de PCData e inicializa os valores
    pcData := &PCData{}
    pcData.Tags.Store(0)
    pcData.UniqueTags.Store(0)
    pcData.CommStatus.Store(false)
    pcData.WifiStatus.Store(false)
    pcData.Lte4Status.Store(false)
    pcData.RfidStatus.Store(false)
    pcData.SysVersion.Store(414)
    pcData.Backups.Store(0)
    pcData.Envios.Store(0)

    // Envia os dados iniciais
    pcData.Send(sender)
    <-time.After(time.Second * 2)

    // Configura um ticker para enviar dados periodicamente
    ticker := time.NewTicker(120 * time.Millisecond)
    defer ticker.Stop()

    log.Println("Iniciando o envio de dados...")
    for range ticker.C {
        pcData.Tags.Add(1)
        pcData.Send(sender)
    }
}
```

## Explicação do Código

1. **Inicialização do SerialSender**: O `SerialSender` é configurado com uma taxa de baud de 115200. Ele gerencia a comunicação serial e reabre a porta automaticamente em caso de falha.

2. **Estrutura PCData**: A estrutura `PCData` armazena informações como tags, status de comunicação, status de Wi-Fi, entre outros. Esses dados são enviados periodicamente através da porta serial.

3. **Envio de Dados**: O método `Send` da estrutura `PCData` formata os dados e os envia usando o `SerialSender`.

4. **Envio Periódico**: Um ticker é usado para enviar os dados a cada 120 milissegundos.

## Licença

Este projeto está licenciado sob a licença MIT. Consulte o arquivo [LICENSE](./LICENSE) para mais informações.

## Contribuição

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

## Referência

Este código foi desenvolvido como parte do post: [Jornada Serial - Parte 1](https://darkcyan-salmon-536746.hostingersite.com/page/jornadaserial1/).
