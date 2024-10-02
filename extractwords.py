import os

words = []
with open("wordFrequency.csv") as csv:
    for line in csv:
        words.append(line.split(",")[1])

if not os.path.exists("wordlists/"):
    os.mkdir("wordlists")

open("wordlists/200.csv", 'a').close()
open("wordlists/1000.csv", 'a').close()
open("wordlists/3000.csv", 'a').close()
open("wordlists/5000.csv", 'a').close()

with open("wordlists/200.csv", "w") as csv200:
    for word in words[1:201]:
        csv200.write(word + "\n")

with open("wordlists/1000.csv", "w") as csv1000:
    for word in words[1:1001]:
        csv1000.write(word + "\n")

with open("wordlists/3000.csv", "w") as csv3000:
    for word in words[1:3001]:
        csv3000.write(word + "\n")

with open("wordlists/5000.csv", "w") as csv5000:
    for word in words[1:5001]:
        csv5000.write(word + "\n")
