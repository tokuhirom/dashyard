# ADR-001: Frontend Framework - React

## Status
Accepted

## Context
We need a frontend framework for the Dashyard dashboard UI. The main candidates considered were React and Svelte.

## Decision
Use React + Vite + TypeScript for the frontend.

## Rationale
- Larger ecosystem with more available libraries and community support
- Better Chart.js integration via react-chartjs-2
- More contributors are likely familiar with React
- Vite provides fast development experience and optimized builds

## Consequences
- Larger bundle size compared to Svelte, but acceptable for a dashboard tool
- Well-established patterns for state management and component composition
