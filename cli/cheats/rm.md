
# Delete the `done` key from the `todo` bucket

slide rm done@todo

# Delete the `browser-history` key from the `topsecret` bucket and purge all history from underlying stream

slide rm browser-history@topsecret --purge
