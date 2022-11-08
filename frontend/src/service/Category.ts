import http from '../http-common'
import Category from '../model/Categories'

const getAll = () => http.get<Array<Category>>("/categories")

const get = (name: string) => http.get<Array<Category>>(`/categories/${name}`)

const post = (category: Category) => http.post<Category>("/categories", category)

const remove = (name: string) => http.delete<string>(`/categories/${name}`)

const CategoryService = {
  getAll,
  get,
  post,
  remove,
};

export default CategoryService;