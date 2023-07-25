# grpc-directory-service
# Описание задачи
Написать асинхронный сервер и клиент, отдающий информацию о директории по заданному в запросе пути (список файлов/директорий и их размер) 
На сервере необходимо поддержать временный кэш (в оперативной памяти) ответов по недавно запрошенным путям. 
Использовать GRPC для общения между сервером и клиентом, реализовать работу с файловой системой встроенными средствами языка.

# Результат работы
![gRPC terminals](https://github.com/Vsevolod-Z/grpc-directory-service/assets/59262675/70f3849b-34c0-457d-9502-36987d47fe37)

# Результаты тестов
![image](https://github.com/Vsevolod-Z/grpc-directory-service/assets/59262675/626d32a9-d474-4c68-bb1c-e63f6f43b4e0)

![image](https://github.com/Vsevolod-Z/grpc-directory-service/assets/59262675/f1ab4d1f-166d-4d0c-b719-1d854a62588c)

![gRPC test](https://github.com/Vsevolod-Z/grpc-directory-service/assets/59262675/81d9a47b-f04a-4762-9718-31da3cf22f6d)
