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

export const Client = axios.create({
  baseURL: 'http://localhost:8080/api/',
  headers: {
    'Content-Type': 'application/json'
  }
})

export const GetAccountByName = async (name: string): Promise<Account> => {
  const response = await Client.get<Account>(`accounts?name=${name}`)
  return response.data
}

export const GetAllAccounts = async (): Promise<Account[]> => {
  const response = await Client.get<Account[]>('accounts')
  const accs = response.data
  return accs
}

export const GetNetWorth = async (): Promise<NetworthPoint[]> => {
  return (await Client.get<NetworthPoint[]>('networth')).data
}

export const GetAllAccountsByClass = async (cls: string): Promise<Account[]> => {
  const response = await Client.get<Account[]>(`accounts?class=${encodeURIComponent(cls)}`)
  return response.data
}
