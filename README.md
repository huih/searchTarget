# searchTarget
the searchTarget is writed in situation that find match string position or no match string position.
searchTarget support three method to search object string position:
1. match method, use set searchType is 2
2. no match method, use set searchType is 1
3. no match up method, use set searchType is 3, this value is default

for example:
target.txt:
test1
test2

checkfile.txt:
test1
sdfdsaf
test1
test2
dfsadfas
test2

use match method : the searchTarget will find next string:
checkfile.txt 3: test1
checkfile.txt 4: test2

the output format is : filename lineno:object string

if you use no match method: the searchTarget will find next string:
checkfile.txt 1: test1
checkfile.txt 2: sdfdsaf

if you use no match up method, the searchTarget will find next string:
checkfile.txt 5: dfsadfas
checkfile.txt 6: test2

