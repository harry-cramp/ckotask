# Take home documentation

## Assumptions

The implementation of this service was carried out with the following assumptions:

- Cards are valid until end of expiry month
- Amount must be greater than zero
- Currency codes are normalised to upper case
- Idempotency not implemented (ran out of time for this)
- Service should fail with `503` if bank service is unavailable
- Rejected payments should not be stored
- Storage is thread-safe (as a real DB would with proper use of transactions)

---

## Architecture

The API structure is split into several layers to aid in readability and maintainability.

- The `domain` layer stores the core business logic such as card verification
- `service` stores abstract methods for interfacing with client layers
- `infrastructure` contains adapters for specific banking services
- `repository` contains methods for interacting with the data store

The `service` layer contains a bank client interface allowing for easy expansion and swapping of specific banking services. This also allows for easier unit testing.

A DDD approach was chosen when designing the validation logic stored in `/domain`. This was done so for various reasons:

- To fail fast, as invalid card data would get rejected at the earliest point
- Invalid data would not propagate through the system at al
- To reduce any duplication of validation logic

However, there is a trade-off as more boilerplate code was required. The approach also caused errors later on in development which made the error handling a bit unorganised when processing the payment.

---

## Testing

The service is covered by several unit and component tests, contained within the folders of the respective code being tested. The aim was to cover the critical paths and business logic, such as failure responses and card validation.

---

## Suggested improvements

With more time, the following changes could have been made:

- Refine the error handling within the `POST /api/payments` handler
- Proper integration tests
- Explicit errors to client (i.e. "amount must be greater than zero")
- Add idempotency checks, either via passing an idempotency key to the bank service (if supported) or by checking the client key against an internal idempotency table with TTL records
- Rate limiting on bank service calls
