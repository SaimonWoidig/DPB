# DPB - 5. Cvičení

Implementujte jednotlivé body pomocí PyMongo knihovny - rozhraní je téměř stejné jako v Mongo shellu.
Před testováním Vašich řešení si nezapomeňte zapnout Mongo v Dockeru.
Funkce find vrací kurzor - pro vypsání výsledku je potřeba pomocí foru iterovat nad kurzorem.
Všechny výsledky limitujte na 10 záznamů. Nepoužívejte české názvy proměnných!

## Základní úlohy
1. Vypsání všech restaurací
2. Vypsání všech restaurací - pouze názvů, abecedně seřazených
3. Vypsání pouze 5 záznamů z předchozího dotazu
4. Zobrazte dalších 10 záznamů
5. #Vypsání restaurací ve čtvrti Bronx (čtvrť = borough)
6. Vypsání restaurací, jejichž název začíná na písmeno M
7. Vypsání restaurací, které mají skóre větší než 80
8. Vypsání restaurací, které mají skóre mezi 80 a 90

## Bonusové úlohy:
9. Vypsání všech restaurací, které mají skóre mezi 80 a 90 a zároveň nevaří americkou (American) kuchyni
10. Vypsání všech restaurací, které mají alespoň osm hodnocení
11. Vypsání všech restaurací, které mají alespoň jedno hodnocení z roku 2014 
12. Uložte novou restauraci (stačí vyplnit název a adresu)
13. Vypište svoji restauraci
14. Aktualizujte svoji restauraci - změňte libovolně název
15. Smažte svoji restauraci
    1. pomocí id (delete_one)
    2. pomocí prvního nebo druhého názvu (delete_many, využití or)

Poslední částí tohoto cvičení je vytvoření jednoduchého indexu.
Použijte např. 3. úlohu s vyhledáváním čtvrtě Bronx. První použijte Váš již vytvořený dotaz a výsledek si vypište na výstup a všimněte si položky 'totalDocsExamined'
Poté vytvořte index na 'borough', zopakujte dotaz a porovnejte hodnoty 'totalDocsExamined'.
