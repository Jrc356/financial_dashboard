import axios from 'axios'

export interface Account {
  Name: string
  Values: AccountValue[]
}

export interface AccountValue {
  Value: number
  Date: string
}

export interface Asset extends Account {
  Type: string
  TaxBucket: string
}

export interface Liability extends Account {}

export enum API {
  Asset = 'asset',
  Liability = 'liability'
}

export const Client = axios.create({
  baseURL: 'http://localhost:8080/api/'
})

export const GetAll = async (api: API): Promise<Account[]> => {
  const response = await Client.get(api)
  const accs = response.data as Account[]
  for (const a of accs) {
    a.Values = []
  }
  return accs
}

export const GetValuesForAccount = async (api: API, account: Account): Promise<Account> => {
  const a = { ...account }
  const response = await Client.get(`${api}/${account.Name}/value`)
  a.Values = response.data as AccountValue[]
  return a
}

export const GetAllAssets = async (): Promise<Asset[]> => {
  const response = await Client.get(API.Asset)
  return response.data as Asset[]
}

export const GetAsset = async (name: string): Promise<Asset> => {
  const response = await Client.get(`${API.Asset}/${name}`)
  return response.data as Asset
}

export const GetAllLiabilities = async (): Promise<Liability[]> => {
  const response = await Client.get(API.Liability)
  return response.data as Liability[]
}

export const GetLiability = async (name: string): Promise<Liability> => {
  const response = await Client.get(`${API.Liability}/${name}`)
  return response.data as Liability
}
