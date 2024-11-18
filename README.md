Aqui está o README em formato Markdown para o seu projeto `RenameFiles`, estilo GitHub, com a lista de parâmetros organizada em uma tabela:

````markdown
# RenameFiles

`RenameFiles` é uma ferramenta de linha de comando para renomear arquivos em um diretório, fornecendo flexibilidade para modificar os nomes dos arquivos de várias maneiras, como substituir caracteres, adicionar timestamps e mudar para título capitalizado.

## Como Usar

Para rodar o comando, utilize:

```sh
renamefiles [flags]
```
````

## Exemplos

- Renomear arquivos no diretório atual substituindo espaços por underscores:

  ```sh
  renamefiles --underscore
  ```

- Renomear arquivos convertendo o nome para Title Case:

  ```sh
  renamefiles --title-case
  ```

- Executar uma simulação (dry-run) sem realmente renomear os arquivos:

  ```sh
  renamefiles --dry-run
  ```

## Parâmetros

| Parâmetro             | Tipo     | Descrição                                                                                     |
| --------------------- | -------- | --------------------------------------------------------------------------------------------- |
| `--underscore`        | Booleano | Substitui espaços por underscores (`_`) nos nomes dos arquivos.                               |
| `--remove-underscore` | Booleano | Substitui underscores (`_`) por espaços nos nomes dos arquivos.                               |
| `--separator`         | String   | Caractere a ser usado como separador (ex.: `_` ou `-`).                                       |
| `--old-separator`     | String   | Caractere a ser substituído nos nomes dos arquivos.                                           |
| `--new-separator`     | String   | Caractere que substituirá o `old-separator` nos nomes dos arquivos.                           |
| `--title-case`        | Booleano | Converte os nomes dos arquivos para Title Case (primeira letra de cada palavra em maiúsculo). |
| `--include-timestamp` | Booleano | Inclui o timestamp de criação do arquivo no nome do arquivo.                                  |
| `--dry-run`           | Booleano | Realiza uma simulação sem renomear os arquivos. Apenas exibe o que seria feito.               |

## Instalação

Para usar o comando `renamefiles`, primeiro você deve instalar a ferramenta. Siga os passos abaixo para compilar e instalar:

```sh
go build -o renamefiles
```

Em seguida, adicione o binário compilado ao seu PATH ou execute diretamente.

## Dependências

Este projeto depende dos seguintes pacotes:

- [Cobra](https://github.com/spf13/cobra) - Para gerenciar os comandos CLI.
- [x/text/unicode/norm](https://pkg.go.dev/golang.org/x/text/unicode/norm) - Para normalização de caracteres Unicode.

## Observações

- Certifique-se de não usar simultaneamente as flags `--underscore` e `--remove-underscore`, pois são opções conflitantes e não funcionarão juntas.
- Quando renomear os arquivos, a ferramenta aplicará todas as opções fornecidas de forma sequencial, afetando o resultado final.

## Licença

Este projeto é distribuído sob a licença MIT. Veja o arquivo [LICENSE](./LICENSE) para mais informações.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir problemas ou enviar pull requests.

## Contato

Se você tiver alguma dúvida ou sugestão, pode entrar em contato comigo através do meu [site pessoal](https://www.robsonalves.dev.br).

```

Este README fornece uma descrição clara sobre o funcionamento do seu programa, os parâmetros disponíveis (organizados em uma tabela) e exemplos de uso. Isso ajuda a garantir que qualquer pessoa que leia o README entenda como utilizar o `RenameFiles` e suas várias funcionalidades.
```
