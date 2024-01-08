import axios from 'axios'

export interface Account {
  name: string
  class: string
  category: string
  taxBucket: string
  values: AccountValue[]
}

export interface AccountValue {
  value: string
  CreatedAt: string
}

export interface NetworthPoint {
  date: Date
  value: number
}

const client = axios.create({
  baseURL: 'http://localhost:8080/api/',
  headers: {
    'Content-Type': 'application/json'
  }
})

export const GetAccountByName = async (name: string): Promise<Account> => {
  const response = await client.get<Account>(`accounts?name=${name}`)
  return response.data
}

export const GetAllAccounts = async (): Promise<Account[]> => {
  const response = await client.get<Account[]>('accounts')
  return response.data
}

export const GetNetWorth = async (): Promise<NetworthPoint[]> => {
  const response = await client.get<NetworthPoint[]>('networth')
  return response.data
}

export const GetAllAccountsByClass = async (cls: string): Promise<Account[]> => {
  const response = await client.get<Account[]>(`accounts?class=${encodeURIComponent(cls)}`)
  return response.data
}
