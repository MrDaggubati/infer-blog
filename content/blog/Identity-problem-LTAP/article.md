---
title: LTAP solves the data movement problem. Who solves the identity problem?
slug: ltap-mdm-identity-problem
description: >
  LTAP brings transactional and analytical workloads closer together, but technical convergence does not solve identity, semantics, relationships, and master data modelling.
date: 2026-10-07
author: Sudhakar Daggubati
tags:
  - MDM
  - Data Modelling
  - LTAP
  - Databricks
  - Snowflake
  - Data Architecture
summary: >
  Poor domain and data modelling is one of the major failure modes I have seen in MDM programs. As Databricks and Snowflake move transactional and analytical workloads closer together, should data platforms also move up the stack and provide reusable identity and master-data foundations? Oracle TCA offers an interesting lesson in how a Party-centric model can prevent every organization from rediscovering the same old MDM pitfalls.

featured: false
---

I once worked on an MDM requirement where I inherited a data model that had already gone through 14 iterations.

End users had been involved. Requirements had been discussed. The model had been reviewed and revised repeatedly.

The project still failed.

The problem was not a lack of tables, attributes, workshops, or iterations. The model itself did not accommodate some of the fundamental needs of Master Data Management.

That experience stayed with me because MDM has an uncomfortable characteristic: a data model can look perfectly reasonable and still be fundamentally wrong.

A `Customer` table can contain every field the business asked for and still be a poor model of a customer.

This becomes particularly interesting now that modern data platforms are moving beyond their analytical origins. Databricks is explicitly developing Lake Transactional/Analytical Processing, or LTAP, around Lakebase. Databricks describes LTAP as an architecture that serves transactional and analytical workloads from a unified storage layer under one governance model.

Snowflake has been moving in a similar direction with Unistore and Hybrid Tables. Hybrid Tables provide row-oriented transactional capabilities, including enforced primary and foreign key constraints, while remaining part of the Snowflake platform and query environment.

This convergence raises a question that I find more interesting than whether the lakehouse can now execute transactions.

**If these platforms are moving closer to operational systems, should they eventually move further up the stack and help customers solve some of the recurring semantic problems around identity and master data?**

Not necessarily by becoming MDM products.

Perhaps by providing some of the reusable foundations that MDM implementations repeatedly have to rediscover.

Oracle's Trading Community Architecture provides an interesting precedent.

## The difficult part of MDM is not creating a Customer table

Consider a simple requirement:

> We need a single view of the customer.

It sounds straightforward.

We might start with something like:

```text
CUSTOMER
--------
customer_id
name
address
email
phone
vat_number
customer_type
```

Then another source appears.

It has another customer ID.

A third source identifies the same organization using a VAT number.

One division treats a subsidiary as a customer. Another treats the parent organization as the customer.

A person can represent an organization.

An organization can be a customer, supplier and partner at the same time.

A customer can have several accounts.

An address can be registered, billing, shipping or temporary.

Two records that look like duplicates might actually represent separate legal entities.

Two records that look different might represent the same entity.

And eventually somebody asks:

> What did we believe about this customer six months ago, and which system told us that?

The original Customer table starts to collapse.

This is where MDM stops being primarily a database design exercise.

It becomes an **identity, semantics, relationship, provenance and governance problem**.

## Customer is not necessarily an identity

One of the most useful ideas in Oracle Trading Community Architecture, TCA, is surprisingly simple.

A party exists independently of the business relationship it has with you.

Oracle's current documentation still describes a party as a person or organization that can enter into business relationships, while the customer represents the selling relationship with that party. A party can have multiple customer accounts representing different business relationships.

Conceptually:

```text
                 PARTY
                   |
          +--------+--------+
          |                 |
       PERSON          ORGANIZATION
          |                 |
          +--------+--------+
                   |
              RELATIONSHIPS
                   |
       +-----------+-----------+
       |           |           |
    Customer    Supplier     Partner
       |
    Accounts
       |
 Account Sites
```

That distinction matters.

"Acme BV" is an organization.

Being my customer is something Acme BV does in relation to my organization.

Acme could simultaneously be my supplier.

John could be a person who is an employee of Acme, a contact for one of its customer accounts, and perhaps independently a customer himself.

The identity should not need to be duplicated simply because its role changes.

Oracle TCA also models relationships explicitly. Oracle describes party relationships in terms of subject, object, role and duration. An employment relationship, for example, relates a person acting as employee to an organization acting as employer.

This gives us a much richer foundation:

```text
IDENTITY
   |
 PARTY
   |
   +---- identifiers
   +---- names and aliases
   +---- locations
   +---- contact points
   +---- classifications
   |
   +---- relationships ---- PARTY
   |
   +---- business roles
            |
            +---- customer
            +---- supplier
            +---- partner
```

It is not a perfect universal model. No model is.

But it captures something important that many MDM projects discover the expensive way:

**Identity and business role are not the same thing.**

## Fourteen iterations can still produce the wrong model

End-user involvement is essential, but it does not guarantee a sound master-data model.

Ask a salesperson whether the Customer screen contains the right information and you may get excellent feedback about names, addresses, segments, contacts and account status.

That does not necessarily answer questions such as:

* What constitutes identity?
* When are two records the same party?
* Can one party have multiple business roles?
* Can several identifiers refer to the same identity?
* Can an identifier move from one party to another?
* Which attributes belong to the party and which belong to the customer relationship?
* How are organizational hierarchies represented?
* Are relationships directional?
* Are relationships time-bound?
* Which source is authoritative for each attribute?
* What happens when authoritative sources disagree?
* How is a merge reversed?
* How is provenance preserved?
* How do we represent what was believed at a particular point in time?

These are not merely user-interface questions.

They are semantic and architectural questions.

A project can therefore iterate repeatedly over the visible model while leaving the underlying abstraction wrong.

## The potholes are remarkably repeatable

This is one reason I consider poor domain and data modelling one of the major failure modes I have seen in MDM programs.

The same traps appear repeatedly.

### Source systems become the enterprise model

CRM has `CUSTOMER`.

ERP has `BUSINESS_PARTNER`.

Billing has `ACCOUNT`.

Support has `CONTACT`.

A common mistake is to combine these structures and call the result a canonical model.

But a canonical model should describe the domain, not merely reconcile the schemas of the systems currently participating in it.

Otherwise today's application landscape becomes tomorrow's enterprise ontology.

### Customer, account and party become interchangeable

This works until one legal entity has several accounts, several legal entities share a commercial arrangement, or the same organization participates in several roles.

Then identity and transactional context become entangled.

### The golden record becomes a golden row

There may be no single authoritative source for an entire party.

CRM might be authoritative for the commercial name.

ERP might own payment terms.

A corporate registry may be authoritative for legal identity.

A verified customer interaction might provide the most recent phone number.

Mastering often occurs at attribute level, not simply at record level.

### History is added afterwards

Operational applications naturally care about the current state.

Master data frequently needs both current state and historical truth.

A relationship can begin and end.

An organization can change name.

A subsidiary can move between parents.

An address can change purpose.

A source can revise a previously supplied value.

Without temporal semantics, the master can tell us what is true now but not necessarily what was true when a business event occurred.

### Relationships are reduced to foreign keys

Real organizations are graphs.

Companies own companies. People work for organizations. Households contain people. Customers have agents. Products belong to multiple classifications.

Trying to represent all of this through a handful of `parent_id` columns eventually becomes restrictive.

Oracle TCA explicitly supports complex and hierarchical party relationships, including relationship direction and date effectivity.

### Matching is treated as an implementation detail

It is not enough to create a `party_id`.

The difficult question is how a source record becomes associated with that party.

```text
CRM-8271  -----+
ERP-99182 -----+----> PARTY-10452
WEB-128   -----+
VAT-NL... -----+
```

Those mappings have lineage.

- Some are deterministic.
- Some are probabilistic.
- Some are manually approved.
- Some should be reversible.

Oracle TCA itself maintains original-system references and includes data-quality machinery around matching and duplicate identification.

Identity resolution is part of the domain architecture, not something to sprinkle on after the tables have been created.

## Now introduce LTAP

This is where the current evolution of Databricks becomes interesting.

The traditional architecture has a fairly obvious separation:

```text
Operational systems
        |
        | CDC / ETL
        v
Analytical platform
        |
        v
BI / ML / AI
```

The operational system owns current business state.

The analytical platform receives a copy.

MDM might sit somewhere between or alongside them.

LTAP challenges part of this separation.

Databricks describes its LTAP architecture as supporting OLTP and OLAP from a unified data storage layer, specifically reducing the need for the CDC, replication and transformation pipelines traditionally required to keep operational and analytical systems synchronized.

Snowflake's Unistore makes a related proposition. Transactional applications can operate on Hybrid Tables while analytical queries can work with transactional and analytical data within the same platform.

That is an important technical change.

But removing the boundary between operational and analytical storage does not remove the semantic problem.

In fact, it may make that problem more visible.

## LTAP can converge storage without converging meaning

Imagine that we successfully remove several copies of customer data.

Excellent.

We now have transactional and analytical workloads operating much closer to the same data.

But we can still have this:

```text
Sales.Customer
Billing.Customer
Support.Customer
Commerce.User
ERP.BusinessPartner
Supplier.Vendor
```

They may all live on the same platform.\
They may all have excellent lineage.\
They may all be protected by the same governance framework.\
They may all be queryable without ETL.\
And they can still disagree about what a customer is.\

That distinction matters:

```text
Technical convergence
        |
        v
shared storage
shared transactions
shared governance
shared lineage
        |
        X
        |
Semantic convergence
        |
        v
shared identity
shared definitions
shared relationships
shared mastering rules
```

**LTAP can solve data movement without solving identity.**

That is not a criticism of LTAP. Identity is a different problem.

But once transactional and analytical workloads begin to converge, it becomes reasonable to ask whether the platform boundary should eventually move upward.

## Unity Catalog governs an object. It does not define the business object

Databricks already has a substantial governance layer in Unity Catalog.

Unity Catalog provides centralized access control, discovery, lineage, auditing, classification and other governance capabilities across data and AI assets.

This is valuable, but there is an important distinction.

A catalog can know:

```text
customer.email
```

It can know who owns it.\
It can know that it contains sensitive information.\
It can know who may read it. \
It can know which pipeline produced it.\
It can know which downstream table consumed it. \

But those capabilities do not inherently answer:

```text
Is this customer a person or an organization?

Is Customer 781 the same party as Supplier 992?

Which identifier establishes the identity?

Is this address a property of the party,
the account,
or a relationship?

Which source wins when two systems disagree?

Was this relationship valid on 12 March 2025?

Can this merge be reversed?
```

Governance metadata and business semantics overlap, but they are not equivalent.

The same applies broadly to Snowflake. Unifying transactional and analytical processing and applying common governance does not automatically create a shared enterprise identity model.

## So should Databricks or Snowflake build an MDM product?

Not necessarily.

That would be too simplistic.

MDM is a large discipline with specialized products, industry-specific requirements, stewardship workflows, data-quality engines, matching algorithms, hierarchy management and governance processes.

There is also a danger in platform vendors becoming too opinionated.

A bank's concept of Party is not identical to a retailer's.\
A healthcare identity model has requirements that a manufacturing supplier model may not.\
B2B organizational hierarchies can be radically different from B2C householding.\

The answer should therefore not be:

> Here is the universal Customer schema. Everyone use it.

That would simply create another rigid canonical model.
There is a more interesting middle ground.

## Build the primitives, not everyone's enterprise model

A modern data platform could provide reusable semantics and capabilities around concepts that recur across industries.

For example:

```text
Party
Party Type
Identifier
Identifier Type
Party Role
Party Relationship
Relationship Type
Location
Contact Point
Source Reference
Attribute Assertion
Effective Period
Provenance
Match Candidate
Match Decision
Merge
Unmerge
Survivorship Rule
```

These could be composable rather than mandatory.

The platform could provide a reference implementation while allowing organizations to extend or replace pieces.

Conceptually:

```text
                    DOMAIN APPLICATIONS
                           |
          +----------------+----------------+
          |                |                |
       Customer         Supplier         Partner
          \                |                /
           +---------------+---------------+
                           |
                 PARTY / IDENTITY LAYER
                           |
          +----------------+----------------+
          |                |                |
      Identity       Relationships      Provenance
          |                |                |
      Matching         Hierarchy       Survivorship
          +----------------+----------------+
                           |
                    GOVERNANCE LAYER
                           |
             policies / lineage / quality
                           |
                     LTAP PLATFORM
                           |
             transactional + analytical
```

That is very different from forcing every customer into a vendor's Customer table.

It is closer to providing a **semantic construction kit for mastered domains**.

## Why platform vendors are in an interesting position to do this

There are several reasons this becomes more plausible as platforms expand toward operational workloads.

### They can see both sides of the lifecycle

Historically, analytical platforms primarily received the result of operational processes.

With transactional capabilities, the platform can increasingly participate in the creation and mutation of operational state.

That creates the possibility of maintaining identity closer to where business state is created rather than reconstructing identity downstream.

### Governance is already there

Databricks already uses Unity Catalog as the governance layer across data and AI.

Snowflake similarly emphasizes common security and governance across transactional and analytical workloads in Unistore.

Identity mastering needs governance.

- Who owns this attribute?
- Where did it originate?
- Who changed it?
- Who may see it?

What downstream assets depend on it?

These platforms already solve substantial parts of that surrounding problem.

### AI makes identity quality more important, not less

There is another reason this matters.

AI systems increasingly consume enterprise context.

An agent asking:

> What is our total exposure to Acme?

has a much harder problem if these are separate identities:

```text
ACME BV
Acme Netherlands
ACME-NL
Customer 8291
Supplier 771
CRM Account 91282
```

Giving an LLM more context does not solve unresolved enterprise identity.

If anything, AI increases the value of a governed semantic identity layer.

Before asking an AI system to reason about a customer, supplier, employee or product, we should have some confidence that the enterprise knows **which real-world entity the data refers to**.

## There is also a strong argument against putting too much into the platform

The opposite argument deserves equal attention.

Enterprise semantics should not become another form of platform lock-in.

If an organization's canonical Party definition exists only because a particular database product understands it, migrating the workloads may become difficult.

This would be especially ironic if open table formats reduced storage lock-in while proprietary semantic models recreated lock-in one level higher.

So any platform-level approach should ideally separate:

```text
Open semantic definition
        |
Reference domain model
        |
Portable contracts
        |
Platform implementation
```

The semantic definition should survive a change in execution engine.

That suggests open schemas, portable contracts, explicit identifiers and APIs may be more important than a proprietary MDM feature branded by a platform vendor.

The vendor can provide an excellent implementation without needing to own the meaning of the enterprise's data.

## Perhaps TCA is useful as a lesson, not a template

I would not suggest copying Oracle TCA into Databricks or Snowflake.

TCA comes from a different application architecture and a different era.

But the architectural lesson remains useful.

Oracle recognized that many applications repeatedly needed to understand the same fundamental things:

- **Who is this?**
- **What kind of party is it?**
- **How do I identify it?**
- **What relationship do I have with it?**
- **How is it related to another party?**
- **Where is it located?**
- **How can it be contacted?**

Those questions have not disappeared because we moved from enterprise databases to cloud data platforms.

If anything, the number of systems providing competing answers has increased.

## We solved copies. Now we need to think about meaning.

The evolution from separate OLTP and OLAP systems toward LTAP is significant.

Databricks is explicitly pursuing an architecture where transactional and analytical workloads can share a governed storage foundation.

Snowflake is pursuing similar convergence through Unistore and Hybrid Tables.

These developments can remove infrastructure, replication and data-movement complexity.

But there is a layer above that technology stack:

```text
Compute
   ↑
Storage
   ↑
Transactions
   ↑
Governance
   ↑
Semantics
   ↑
Identity
   ↑
Business domain
```

The industry has invested enormously in the lower layers.

MDM programs repeatedly remind us how difficult the upper layers remain.

Poor domain and data modelling is one of the major reasons I have seen MDM initiatives struggle or fail. The potholes are numerous, and many organizations discover the same ones independently.

- Party versus customer.
- Identity versus account.
- Entity versus role.
- Current state versus historical truth.
- Record lineage versus attribute provenance.
- Duplicate detection versus identity resolution.
- Golden record versus survivorship.
- Hierarchy versus relationship.
- These are not new problems.

Oracle's TCA addressed several of them decades ago through an opinionated Party model. Modern architectures should not blindly reproduce TCA, but there is little value in thousands of organizations independently falling into the same old potholes.

As Databricks, Snowflake and other data platforms move closer to operational workloads, perhaps the next interesting step is not simply more database capability.

Perhaps it is providing reusable, open and extensible foundations for **identity, relationships and mastered domains**, while leaving organizations in control of their own business semantics.

LTAP can bring the transaction and the analysis together.

The harder question is whether the industry can also get better at agreeing on **who or what that transaction was actually about**.