// Copyright 2022 MrTimeout
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := execute(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.Default())

	getInstance().AutoMigrate(new(category), new(word))

	main := r.Group("/api")

	categories := main.Group(CategoriesPath)
	{
		categories.GET("", getCategories)
		categories.POST("", addCategory)

		categories.GET(CategoryByIdPath, getCategoryByName)
		categories.DELETE(CategoryByIdPath, delCategory)
	}

	words := categories.Group(WordsPath)
	{
		words.GET("", getWords)
	}

	r.Run(":9090")
}
