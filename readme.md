# Gossip мультилог



## Как развернуть локально два инстанса:
1. `go build`
2. Сделать две отдельных папки для инстансов:
```bash
mkdir send recv
cp -r go-log-gossip send
cp -r go-log-gossip recv
```
3. Внутри каждой папки сделать `./go-log-gossip i`
4. Поменять конфиг внутри каждой ноды для отправки и получения:
   1. Примерный конфиг для отправки:
    ```yaml
      paths:
         path_to_folder_with_rsa_keys: .keys
         path_to_folder_with_logs: logs
      gossip:
         self_node_name: denis-sender
         self_node_port: 5555
         secret_key: Y2hhbmdlQUVTa2V5MDAwMA==
         boostrap_node_addr: 127.0.0.1:8888
         is_boostrap_node: false
      api:
       addr: :7777
    ```
   2. Примерный конфиг для получения:
   ```yaml
    paths:
       path_to_folder_with_rsa_keys: .keys
       path_to_folder_with_logs: logs
    gossip:
       self_node_name: denis-recv
       self_node_port: 8888
       secret_key: Y2hhbmdlQUVTa2V5MDAwMA==
       boostrap_node_addr: 127.0.0.1:8888
       is_boostrap_node: true
    api:
       addr: :5000
   ```
5. Запустить recv: `cd recv/ && ./go-log-gossip s`
6. Запустить send: `cd send/ && ./go-log-gossip s`


Жданов Илья Александрович ФТ-303
Сурков Денис Дмитриевич ФТ-301