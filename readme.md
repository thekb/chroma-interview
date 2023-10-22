# Key Value Store with Transactions

Implement a Key value store with transactions in a language of your choice.
Please implement the following high level API for a key value store along with tests.

## Initialization
```
store = new Store(<anything>)
```
## Data Manipulation
```
store.set(K, V)
store.get(K)
store.delete(K)
Transactions
store.begin() # begins the transaction
store.commit() # commits the transaction 
store.rollback() # rolls the transaction back
```
Feel free to implement additional public methods to this API that you think would be useful for a user of this store.

## Example
```
store.begin()
store.set(“a”, 50)
store.begin()
store.set(“a”, 60)
store.rollback()
store.commit() 
expect store.get(“a”) == 50
```
