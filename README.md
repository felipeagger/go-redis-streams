# go-redis-streams
Demo using messaging with Redis Streams in Golang

Video: https://youtu.be/wkTmbf-WnU8

# Redis

**O que e?**

Redis e um Banco de dados não relacional OpenSource, que tem dentro de sua estrutura o armazenamento chave-valor.
O Redis tem estratégias para guardar os dados em memória e em disco, garantindo resposta rápida e persistência de dados. Os principais casos de uso do Redis incluem cache, gerenciamento de sessões, PUB/SUB.

# Redis Streams para Mensageria (ou Messaging)

![Design of flow](/media/flow.png)

**Pontos Positivos**

- Suporta Topicos e Filas 
- Persistencia em disco (através dos arquivos RDB)
- Alta disponibilidade (com Clusterizacao)
- Alto Throughput
- Permite Reprocessamento
- Possui Consumer Groups
- Latencia minima
- Nao necessita de zookeper
- Ocupa muito menos recursos em relacao ao (Kafka/RabbitMQ)

**Pontos Negativos**

- Nao garante ordem de entrega (ainda)
- Msgs processadas com error nao retorna para redistribuicao


# Links

https://www.youtube.com/watch?v=JpeHIbzmGP4

https://redis.io/topics/streams-intro

https://redislabs.com/blog/use-redis-streams-apps/

https://redislabs.com/blog/getting-started-with-redis-streams-and-java/

