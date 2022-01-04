## Overview

## Summary of NFTContract

NFTContract is a Non Fungible Token (NFT) standard for Flow blockchain.
It offers a powerful set while keeping unnecessary complexity to a minimum and focus on efficiency.
Our Contract consists of different parts like Brand, Schema, Templates and NFTs, which serve different purposes and often reference each other.

## Brand

A brand is an intangible marketing or business concept that helps people identify a company, product, or individual. People often confuse brands with things like logos, slogans, or other recognizable marks. Brand is a symbol of organization that represent parent company. Organizations can create their own Brands in FLOW Blockchain using NFTContract. We have set of restrictions on flow Accounts. In order to create a brand you need to be whitelisted(Approval of Super Admin).

### Schema

Schemas are used to define a data structure of NFT. Schema is like an interface that used for creating a new template(defined later) using the structure we defined in schema. A Schema is a collection of Key attributes. A schema is owned by an account. Schema objects are logical structures. Keys and Values in schemas hold data types, or can consist of a definition only, such as a view or synonym. We are supporting following datatypes in schema:

- String
- Int
- Fix64
- Bool
- Address
- Array
- Any

## Template

Templates are blueprints of NFTs. For creating NFTs, we use Templates as defined schemas. Flow Blockchains are storing metadata offchain but only we are creating a structure where we can store metadata onchain.
