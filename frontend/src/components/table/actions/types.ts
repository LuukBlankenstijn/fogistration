import type { Table } from "@tanstack/react-table"

export interface HeaderInput<T> {
  rows: T[],
  table: Table<T>
}

export interface RowInput<T> {
  row: T,
  table: Table<T>
}

export interface HeaderAction<T> {
  label: string
  onClick: (data: HeaderInput<T>) => void | Promise<void>
  disabled?: (data: HeaderInput<T>) => boolean
}

export interface RowAction<T> {
  label: string
  onClick: (data: RowInput<T>) => void | Promise<void>
  disabled?: (data: RowInput<T>) => boolean
}
