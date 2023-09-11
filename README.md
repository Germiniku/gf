# gt
gorm kits help user write gorm curd happy!

when we `dao` incoming structure field has tag, tag auto stitch query conditions

| Opr           | Description |
|---------------|-------------|
| eq            | equal       |
| neq           | not equal   |
| gt            | gt          |
| gte           | gte         |
| lt            | lt          |
| lte           | lte         |
| in            | IN          |
| !in           | NOT IN      |
| like/contains | LIKE        |
| !like/!contains       | NOT LIKE    |
| any       | = ANY()     |
| overlap       | && ARRAY[]  |

## best practices
```go
package dao

import (
	"context"
	"gorm.io/gorm"
	"github.com/Germiniku/gf"
)

type Dao struct {
	db *gorm.DB
}

type Book struct {
	Name string `filter:"col:name;opr:eq"`
}

func (d *Dao) Books(ctx context.Context, q *Book) (books []*Book,err error){
  books = make([]*Book,0)
  err = d.db.Model(&Book{}).Scopes(gf.Filter(q)).Find(&books).Error
  return 
}
```