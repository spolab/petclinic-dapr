# Broadcaster

The broadcaster is a microservice that reads events out of the MongoDB event log and publishes them as Cloud Events.

## What problem does it solve

In event sourcing, you need to make sure that the events and the snapshot are written in the same transaction.

Let's see what are the consequences if you do not:

### Case 1 - update the aggregate state and then publish the events

In this case, there is a risk that the snapshots are up-to-date, but the materialized view will not be updated because of an error while broadcasting the events. Commands will not be able to interact with that aggregate ID because their version will always be one version late.

### Case 2 - publish the events first and then the snapshots

In this case, there are two risks:

- The events get propagated but then there is an error while updating the aggreate state. This problem can potentially create a data integrity as the aggregate's state and the materialized views will be out of sync. Any further command will additionally extend the integrity problem and, 
- The reader may observe that, in this case, all one needs to do is to replay the event log. While this is certainly correct, there may still be a quantum of time when the materialized views are up-to-date and the aggregate state is not.

## How does it solve it

Aggregate roots in this demo update their state and append the new events *in the same database transaction*. This approach provides two advantages:

- Ensures consistency between state and events: if the state cannot be updated then the events will not be published, and vice versa.
- Leverages the locking mechanism of the database, so that other instances of the same aggregate will not execute out-of-order updates. 

So, the broadcaster listens to the database changes 