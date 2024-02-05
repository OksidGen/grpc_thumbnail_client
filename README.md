# gRPC Thumbnail Client
gRPC Thumbnail Client - это клиент, написанный на языке программирования Go, для работы с [gRPC Thumbnail Server](https://github.com/OksidGen/grpc_thumbnail_server.git).

## Функциональность
- Запрос превью с сервера по ссылке на видео
- Конвертирование полученных байтов в изображение и сохранение

## Установка и запуск
Склонируйте репозиторий:
```
git clone https://github.com/OksidGen/grpc_thumbnail_client.git
cd grpc_thumbnail_client
```
Зависимости для проекта указаны в файле go.mod. Вы можете установить их с помощью команды:
```
go mod download
```
Сгенерируйте файлы .pb.go из файла .proto при помощи команды:
```
protoc --proto_path=./proto --go_out=. --go-grpc_out=. ./proto/thumbnail_service.proto
```
Скомпилируйте приложение:
```
go build cmd/grpc_thumbnail_client.go
```
Запустите приложение и передайте необходимые параметры
```
.\grpc_thumbnail_client.exe --async --output_dir=./thumbnails link1 link2
```

## Передаваемые параметры

--async - конфигурирует приложение для синхронного/асинхронного режима
--output_dir - путь до папки для сохранения превью
link1,link2 - ссылки на видеоролики YouTube

## Планы на будущее
 
- [ ] Покрытие кода тестами (в процессе 🚀)
- [ ] Добавление возможности работы по файлу .txt
