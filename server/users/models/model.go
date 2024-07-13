package models

import "github.com/satyambaran/oauth/server/users/structs"

var Models []interface{} = []interface{}{
    &structs.User{},
    &structs.Resource{},
}
