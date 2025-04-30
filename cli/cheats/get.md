
# Gets a key named `test` from the default bucket

slide get test

# Gets a key named `verysecret` in the `secret` bucket with value `hello` and client side decryption using private key seed

slide get verysecret@secret --client-encryption

# Gets a key named `verysecret` in the `secret` bucket with value `hello` and client side decryption using a provided curve key seed

slide get verysecret@secret --client-encryption --seed SXADDZVOUPJO7GVIVYF3JGPODPMLS2RB2J3UZUVO6ONXKVN6NEIZPRYPBQ

# Gets a key named `todo` in the `things` bucket and writes contents to the file `todo.txt`

slide get todo@things > todo.txt

# Gets a key named `logo` in the `images` bucket, base64 decodes and write contents to the file `logo.png`

slide get logo@images | base64 -d > logo.png
