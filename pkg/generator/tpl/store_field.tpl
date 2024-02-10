		if req.{{.Name}} != nil {
    		db.Where("{{.NameSnake}} = ?", req.{{.Name}})
		}