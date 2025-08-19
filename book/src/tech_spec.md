# Technical Specification

## Purpose of this Document

This document defines the technical architecture, design, and functionalities of Khedra, enabling developers and engineers to understand its internal workings and design principles. For a less technical overview of the application, refer to the [User Manual](../user_manual/).

## Intended Audience

This specification is for:

- Developers working on Khedra or integrating it into applications.
- System architects designing systems that use Khedra.
- Technical professionals looking for a detailed understanding of the system.

## Scope and Objectives

The specification covers:

- High-level architecture.
- Core functionalities such as blockchain indexing, REST API, and address monitoring.
- Design principles, including scalability, error handling, and integration with IPFS.
- Supported chains, RPC requirements, and testing methodologies.

## System Overview

Khedra is a sophisticated blockchain indexing and monitoring solution designed with a local-first architecture. It creates and maintains the Unchained Index - a permissionless index of address appearances across blockchain data - enabling powerful monitoring capabilities for any address on any supported EVM-compatible chain.

### Core Technical Components

1. **Indexing Engine**: Processes blockchain data to extract and store address appearances
2. **Service Framework**: Manages the lifecycle of modular services (scraper, monitor, API, IPFS, control)
3. **Data Storage Layer**: Organizes and persists index data and caches
4. **Configuration System**: Manages user preferences and system settings
5. **API Layer**: Provides programmatic access to indexed data

### Key Design Principles

Khedra's technical design adheres to several foundational principles:

1. **Local-First Processing**: All data processing happens on the user's machine, maximizing privacy
2. **Chain Agnosticism**: Support for any EVM-compatible blockchain with minimal configuration
3. **Modularity**: Clean separation of concerns between services for flexibility and maintainability
4. **Resource Efficiency**: Careful management of system resources, especially during indexing
5. **Resilience**: Robust error handling and recovery mechanisms
6. **Extensibility**: Interfaces intended to allow additional components without refactoring core code

## Technology Stack

Khedra is built on a modern technology stack:

- **Go**: The primary implementation language, chosen for its performance, concurrency model, and cross-platform support
- **IPFS**: For distributed sharing of index data
- **RESTful API**: For service integration and data access
- **YAML**: For configuration management
- **Structured Logging**: For operational monitoring and debugging

## Target Audience

This technical specification is intended for:

- **Developers**: Contributing to Khedra or building on top of it
- **System Administrators**: Deploying and maintaining Khedra instances
- **Technical Architects**: Evaluating Khedra for integration with other systems
- **Advanced Users**: Seeking a deeper understanding of how Khedra works

## Document Structure

The remaining sections of this specification are organized as follows:

- **System Architecture**: The high-level structure and components
- **Core Functionalities**: Detailed explanations of key features
- **Technical Design**: Implementation details and design patterns
- **Supported Chains**: Technical requirements and integration details
- **Command-Line Interface**: API and usage patterns
- **Performance and Scalability**: Benchmarks and optimization strategies
- **Integration Points**: APIs and interfaces for external systems
- **Testing and Validation**: Approaches to quality assurance
- **Appendices**: Technical reference materials

This specification aims to provide a comprehensive understanding of Khedra's technical aspects while serving as a reference for implementation, maintenance, and extension of the system.
