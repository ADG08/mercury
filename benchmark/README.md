# Résultats des Benchmarks k6

## Conditions d'exécution

- **Machine** : AMD Ryzen 5 5600X 6-Core Processor (12 processeurs logiques), 16 Go de RAM.
- **Commande k6 exacte (PowerShell)** :

	```powershell
	cat .\benchmark\k6\read-heavy.js | docker run --rm -i grafana/k6 run -
	```

- **État avant chaque tir** :
	- **In-Memory** : redémarrer le serveur afin de repartir avec une map vide ; aucune base PostgreSQL ni aucun cache Redis utilisé.
	- **PostgreSQL** : exécuter `TRUNCATE TABLE urls;` avant le démarrage du serveur.
	- **Redis Cache** : exécuter `redis-cli FLUSHALL` avant le démarrage du serveur.

Scénario : `read-heavy.js` (150 VUs max, 50s, 5 codes distincts avec popularité pondérée, 99.9% lectures `GET /{code}`)

| Stack / Variante | Débit (req/s) | Latence moy. | p(95) | p(99) | Taux d'erreur |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Go (In-Memory)** | 1 542,0 | 1,31 ms | 2,57 ms | 4,12 ms | 0,00 % |
| Go (PostgreSQL) | - | - | - | - | - |
| Go (Redis Cache) | - | - | - | - | - |
| .NET (In-Memory) | - | - | - | - | - |
| .NET (PostgreSQL) | - | - | - | - | - |
| .NET (Redis Cache) | - | - | - | - | - |
| Java (In-Memory) | - | - | - | - | - |
| Java (PostgreSQL) | - | - | - | - | - |
| Java (Redis Cache) | - | - | - | - | - |
| Python (In-Memory)| - | - | - | - | - |
| Python (PostgreSQL)| - | - | - | - | - |
| Python (Redis Cache)| - | - | - | - | - |