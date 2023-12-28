import axios from 'axios'

export interface Account {
  name: string
  class: string
  category: string
  taxBucket: string
  values: AccountValue[]
}

export interface AccountValue {
  value: number
  CreatedAt: string
}

export const Client = axios.create({
  baseURL: 'http://localhost:8080/api/'
})

export const GetAccountByName = async (name: string): Promise<Account> => {
  const response = await Client.get(`accounts?name=${name}`)
  return response.data as Account
}

export const GetAllAccounts = async (): Promise<Account[]> => {
  const response = await Client.get('accounts')
  const accs = response.data as Account[]
  for (const a of accs) {
    a.values = []
  }
  return accs
}

export const GetValuesForAccount = async (account: Account): Promise<Account> => {
  const a = { ...account }
  const response = await Client.get(`accounts/value?name=${encodeURIComponent(account.name)}`)
  a.values = response.data as AccountValue[]
  return a
}

export const GetAllAccountsByClass = async (cls: string): Promise<Account[]> => {
  const response = await Client.get(`accounts?class=${encodeURIComponent(cls)}`)
  return response.data as Account[]
}
