## Queries
 
Pequeña librería para hacer paginar y hacer preloading dinámico, 
hace uso de echo web framework y GORM.

Note: _puede usarse con otros framewors_


```go
    type User struct {
        ID       uint         `json:"id"`
        Nombre   string       `json:"nombre"`
        Permisos []Permiso    `json:"permisos"`
        Targetas []CreditCard `json:"targetas_credito"`
    }
    // Implementación de la interaface Preloader
    // retornamos los campos a los que podemos hacer preloading
    // estos campos son el nombre puesto en el json tag
    // este valor será traducido de snackecase a camelcase
    // por tanto, si estos nombres no coinciden con el atributo de la
    // estructura, podemos especificar explícitamente el nombre
    // pasando como segundo parámetro, separado con una coma
    // como es el caso de `"targetas_credito,Targetas"`
    func (*User) Preload()[]string{        
        return []string{"permisos","targetas_credito,Targetas"}
    }

    type Permiso struct {
        ID   uint   `json:"id"`
        Perm string `json:"perm"`
    }
    type CreditCard struct {
        ID     uint   `josn:"id"`
        Numbre string `json:"content"`
    }
```

En echo debemos agregar el middleware `QueryParamMiddl`, este agregará los campos query parama en el contexto de la petición http

```go
    // Model recibe el contexto y un modelo que implemente la interface
    // Preloader. Es necesario solo si queremos hacer preloading
    ctx=queries.Model(ctx, &models.User{}) 
```

```go
    // hacemos wrapp a nuestra conexión *gorm.DB antes de
    // realizar la consulta.
    conn := queries.Wrapp(ctx, gormConnection)

    //tx := queries.Wrapp(ctx, database.Conn(ctx))
```

     Finalemente podemos hacer consultas:
    http://localhost.com/usuarios?preload=permisos,targetas_credito&limit=100&offset=200