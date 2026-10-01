# Mercury

Un petit laboratoire perso pour construire un URL shortener, le charger avec k6 et comparer des choix d'architecture avec des chiffres réels.

L'idée est simple : partir de l'in-memory, ajouter PostgreSQL puis Redis, et voir concrètement ce que chaque étape change en débit, latence et erreurs. Pas de benchmark théorique : même API, même scénario, conditions documentées.

J'ai choisi ces stacks par curiosité : **Go** est devenu populaire et j'ai envie de voir ce qu'il donne vraiment en performances. **Java** et **.NET**, on en voit dans beaucoup d'offres et de boîtes, souvent derrière des applications qui tournent depuis longtemps. Et **Python**, même si j'accroche moins avec le langage, j'ai quand même envie de voir ce que donnent ses chiffres.

## État actuel

La première version tourne en Go avec `net/http` et un stockage en mémoire.

```powershell
go run .\apps\go\main.go
```

API : `POST /url` pour créer un raccourci, puis `GET /{code}` pour obtenir une redirection `307`.

Principe volontaire : chaque `POST /url` crée un nouveau raccourci. Même si la même URL existe déjà, elle n'est pas dédupliquée et crée quand même une nouvelle ligne/entrée.

Baseline actuelle : **1 542 req/s**, latence moyenne **1,31 ms**, p(95) **2,57 ms**, p(99) **4,12 ms**, **0,00 %** d'erreurs.

## Benchmark

```powershell
cat .\benchmark\k6\read-heavy.js | docker run --rm -i grafana/k6 run -
```

Le scénario monte à 150 VUs, crée quelques URLs distinctes au démarrage puis lit leurs codes avec une popularité volontairement inégale. Les conditions, la machine et les résultats détaillés sont dans [benchmark/README.md](benchmark/README.md).

## Conclusion

À compléter avec mes étonnements, les difficultés d'implémentation, les concepts intéressants de chaque langage et les sources utiles associées.

