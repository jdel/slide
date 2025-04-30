
# Creates a key named `test` in the default bucket with value `hello`

slide set test hello

# Creates a key named `verysecret` in the `secret` bucket with value `hello` and client side encryption using private key seed

slide set verysecret@secret hello --encrypted

# Creates a key named `verysecret` in the `secret` bucket with value `hello` and client side encryption using a provided curve key seed

slide set verysecret@secret hello --encrypted --seed SXADDZVOUPJO7GVIVYF3JGPODPMLS2RB2J3UZUVO6ONXKVN6NEIZPRYPBQ

# Create a key named `todo` in the `things` bucket with contents of the file `todo.txt`

slide set todo@things < todo.txt

# Create a key named `logo` in the `images` bucket with contents of the file `logo.png` base64 encoded

slide set logo@images $(cat logo.png | base64)
