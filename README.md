Fixes grammar in text one sentence at a time.
for all files in a directory structure, passing ""path" param as the root folder.

ask chat gpt to correct any grammar punction errors in each sentence
since there are different responses and some false positives run each sententence up to 8 times.
break out of the loop if it claims no errors.
if it never says that there's a good chance it found an error that is not a false positive.
which ever correction has the most repeated answers from the questions is typically correct, so use that one.
reassemble sentenes into paragraphs
overwrite the existing file with the updated data .

params:
    -apikey <apikey for chatgpt>
    -path <folder containing text files>



building...

go build -o fixer.exe ./src