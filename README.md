# Ergos

### Asynchronous CloudEvents Choreography & Concurrency Ingress Network for Agentic RWA Tokenization

Ergos is an event coordination layer designed to connect asynchronous agent workflows with **Drunix** for Real-World Asset (RWA) tokenization.

Autonomous agents may independently produce valuation, compliance, treasury, and other signals in different formats and at different times. Instead of sending every partial or overlapping request directly to the ledger, Ergos normalizes and correlates these events, waits for the required signals, creates a deterministic transaction intent, and sequences transactions by asset before submitting them to Drunix.

**Drunix remains the authoritative system of record** and is responsible for smart-contract execution, endorsement, final validation, and ledger state.

---

## Problem

Agentic financial workflows are inherently asynchronous.

For example, an RWA transaction may require:

- Asset valuation approval
- Compliance verification
- Treasury/funds verification

These signals can arrive independently and concurrently.

Directly forwarding every agent request to the ledger can create unnecessary contention when multiple transactions target the same underlying asset or ledger state.

Ergos provides a coordination layer between these autonomous agents and Drunix.

---

## Architecture

```text
                  AGENT FLEET
        ┌────────────┬────────────┐
        │ Valuation  │ Compliance │ Treasury
        │   Agent    │    Agent   │  Agent
        └─────┬──────┴──────┬─────┘
              │             │
              │ CloudEvents │
              ▼             ▼
        ┌─────────────────────────┐
        │         Ergos         │
        │                         │
        │  CloudEvent Ingress     │
        │          ↓              │
        │  Event Correlation      │
        │          ↓              │
        │  Correlation Barrier    │
        │          ↓              │
        │  Transaction Intent     │
        │          ↓              │
        │  Asset-aware Sequencer  │
        └────────────┬────────────┘
                     │
                     │ Deterministic
                     │ Transaction Intent
                     ▼
        ┌─────────────────────────┐
        │         DRUNIX           │
        │                         │
        │  Gateway / Endorsement  │
        │          ↓              │
        │  Smart Contract         │
        │          ↓              │
        │  MVCC Validation        │
        │          ↓              │
        │  Authoritative Ledger   │
        └─────────────────────────┘
````

### Responsibility split

| Component     | Responsibility                                        |
| ------------- | ----------------------------------------------------- |
| Agent Fleet   | Produce independent asynchronous signals              |
| Ergos       | Parse, correlate, coordinate and sequence events      |
| Drunix        | Execute transactions and maintain authoritative state |
| RWA Chaincode | Enforce valid asset state transitions                 |

Ergos does **not** replace Drunix's transaction validation or ledger.

---

## Prototype Flow

The current prototype demonstrates an RWA asset moving from:

```text
AVAILABLE → LOCKED
```

### 1. Agent events arrive

Three CloudEvents are produced for the same transaction:

```text
asset.valuation.completed
asset.compliance.completed
asset.funds.verified
```

All events contain:

```text
Correlation ID: txn-001
Asset ID:       SOLAR-001
```

### 2. Ergos correlates the events

Ergos stores events by correlation ID and waits until all required event types are present.

```text
Valuation  ──┐
Compliance ──┼──> Correlation Barrier
Treasury   ──┘
                  │
                  ▼
             Barrier Ready
```

### 3. A deterministic transaction intent is created

Once the barrier is satisfied:

```text
Correlation ID: txn-001
Asset ID:       SOLAR-001
Action:         LOCK_ASSET
```

### 4. Asset-aware sequencing

The transaction is placed into a queue associated with the affected asset.

This allows transactions affecting the same asset/state key to be ordered while keeping independent asset queues separate.

The goal is to **reduce avoidable contention**, not to bypass Drunix's validation.

### 5. Ergos submits the transaction to Drunix

Drunix receives the deterministic transaction intent and executes the RWA chaincode:

```text
LockAsset("SOLAR-001")
```

### 6. Drunix becomes the source of truth

The final state is read back from Drunix:

```json
{
  "id": "SOLAR-001",
  "owner": "issuer-001",
  "total_tokens": 100000,
  "available_tokens": 100000,
  "status": "LOCKED",
  "settlement_status": "PENDING"
}
```

---

## What the Prototype Demonstrates

* [x] CloudEvents ingestion and parsing
* [x] Agent event normalization
* [x] Correlation using transaction IDs
* [x] Multi-agent correlation barrier
* [x] Deterministic transaction intent generation
* [x] Asset-aware transaction sequencing
* [x] Custom RWA smart contract on Drunix
* [x] Real transaction submission to Drunix
* [x] Authoritative ledger state transition
* [x] Reading the resulting state back from Drunix

This is intentionally a **small architectural proof-of-concept**, rather than a production deployment.

---

## Example

Three agents independently produce:

```text
Valuation Agent
    ↓
asset.valuation.completed

Compliance Agent
    ↓
asset.compliance.completed

Treasury Agent
    ↓
asset.funds.verified
```

Ergos correlates them:

```text
txn-001
SOLAR-001
    │
    ▼
All required signals received
    │
    ▼
LOCK_ASSET
    │
    ▼
Asset-aware Sequencer
    │
    ▼
Drunix
    │
    ▼
SOLAR-001 = LOCKED
```

---

## Project Structure

```text
arealis-drunix-hackathon/
│
├── architecture/
│
├── chaincode/
│   └── rwa/
│       ├── asset.go
│       ├── contract.go
│       ├── main.go
│       └── go.mod
│
├── cmd/
│   └── demo/
│       └── main.go
│
├── drunix/
│   └── client.go
│
├── examples/
│   └── solar-asset/
│       └── events.json
│
├── handler/
│   ├── ingress.go
│   └── ingress_test.go
│
├── models/
│   ├── cloud_event.go
│   ├── agent_event.go
│   └── transaction.go
│
├── repository/
│   └── correlation_store.go
│
├── schemas/
│   └── event_types.go
│
└── services/
    ├── correlation/
    │   ├── barrier.go
    │   └── barrier_test.go
    │
    ├── orchestration/
    │   ├── service.go
    │   └── service_test.go
    │
    └── sequencer/
        ├── sequencer.go
        └── sequencer_test.go
```

---

## Running the Prototype

### Prerequisites

* Go 1.25+
* Docker
* Running Drunix test network
* Drunix channel and RWA chaincode deployed

### Start the Drunix network

From the Drunix test-network directory:

```bash
./network.sh up createChannel -c mychannel
```

### Deploy the RWA chaincode

From the Drunix test-network directory:

```bash
./network.sh deployCC \
  -ccn rwa \
  -ccp ../../../arealis-drunix-hackathon/chaincode/rwa \
  -ccl go
```

### Run the Ergos demo

From the project root:

```bash
go run ./cmd/demo
```

Expected flow:

```text
Starting Ergos demo...

Connecting to Drunix...
Connected to Drunix successfully!

Creating RWA asset...
RWA asset created!

Received 3 agent events

Agent event: asset.valuation.completed
  → Correlation barrier not ready

Agent event: asset.compliance.completed
  → Correlation barrier not ready

Agent event: asset.funds.verified
  → Barrier satisfied
  → Queued intent: LOCK_ASSET SOLAR-001

Sequencer produced intent:
  Action: LOCK_ASSET

Submitting transaction to Drunix...

Transaction committed successfully!

Final Drunix ledger state:
SOLAR-001 → LOCKED
```

---

## Design Principle

The key design principle is:

> **Ergos coordinates agent activity; Drunix owns financial state.**

Ergos is responsible for handling the asynchronous nature of agent workflows and producing cleaner, deterministic transaction intents.

Drunix remains responsible for:

* Smart-contract execution
* Transaction endorsement
* MVCC validation
* Ledger persistence
* Authoritative RWA state

Ergos's sequencing is therefore a **pre-coordination mechanism**, not a replacement for Drunix's concurrency control.

---

## Current Prototype Scope

This repository intentionally demonstrates only the core architectural path:

```text
Agent Events
     ↓
Correlation
     ↓
Transaction Intent
     ↓
Sequencing
     ↓
Drunix
     ↓
RWA State Transition
```

The prototype does not attempt to implement a production-grade distributed event platform.

---

## Future Hackathon Work

Potential extensions during the hackathon include:

* Real agent ingress instead of static example events
* Additional RWA lifecycle transitions
* Transfer and settlement workflows
* Conflict-aware scheduling and benchmarking
* Event persistence and TTL handling
* Retry and failure recovery
* Higher-throughput event processing
* Integration with private NPCI/Citi infrastructure
* Production deployment architecture
* Observability and operational tooling

The effectiveness of conflict-aware sequencing would be evaluated through benchmarks rather than assuming that it eliminates MVCC conflicts.

---

## Status

**Prototype: Working**

The current implementation demonstrates an end-to-end flow where asynchronous agent events are coordinated by Ergos and result in an actual authoritative RWA state transition on Drunix.

```text
CloudEvents
     ↓
Ergos
     ↓
Transaction Intent
     ↓
Sequencer
     ↓
Drunix
     ↓
AVAILABLE → LOCKED
```
