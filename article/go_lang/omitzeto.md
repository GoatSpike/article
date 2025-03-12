

omitemptyあり
```
type RippleDest struct {
	UserID              int64   `json:"userID"`
	UserNo              string  `json:"userNo"`
	UserName            string  `json:"userName"`
	UserRomanName       *string `json:"userRomanName,omitempty"`
	UserImage           *string `json:"image,omitempty"`
	DepartmentID        int64   `json:"departmentID"`
	DepartmentName      string  `json:"departmentName"`
	DepartmentRomanName *string `json:"departmentRomanName.omitempty"`
	BaseID              int64   `json:"baseID"`
	BaseName            string  `json:"baseName"`
	TypeDiv             int     `json:"typeDiv"`
}
```

```
 {
    "baseID": 6,
    "baseName": "GBS",
    "departmentID": 5,
    "departmentName": "SATO Singapore Group 1",
    "typeDiv": 2,
    "userID": 324,
    "userName": "（TLB）TEST USER 05",
    "userNo": "555",
    "userRomanName": ""
},

何もなし
```
type RippleDest struct {
	UserID              int64   `json:"userID"`
	UserNo              string  `json:"userNo"`
	UserName            string  `json:"userName"`
	UserRomanName       *string `json:"userRomanName"`
	UserImage           *string `json:"image"`
	DepartmentID        int64   `json:"departmentID"`
	DepartmentName      string  `json:"departmentName"`
	DepartmentRomanName *string `json:"departmentRomanName"`
	BaseID              int64   `json:"baseID"`
	BaseName            string  `json:"baseName"`
	TypeDiv             int     `json:"typeDiv"`
}
```

```
{
    "baseID": 6,
    "baseName": "GBS",
    "departmentID": 5,
    "departmentName": "SATO Singapore Group 1",
    "departmentRomanName": null,
    "image": null,
    "typeDiv": 2,
    "userID": 324,
    "userName": "（TLB）TEST USER 05",
    "userNo": "555",
    "userRomanName": ""
},
```



```

type RippleDest struct {
	UserID              int64   `json:"userID"`
	UserNo              string  `json:"userNo"`
	UserName            string  `json:"userName"`
	UserRomanName       *string `json:"userRomanName,omitzero"`
	UserImage           *string `json:"image,omitzero"`
	DepartmentID        int64   `json:"departmentID"`
	DepartmentName      string  `json:"departmentName"`
	DepartmentRomanName *string `json:"departmentRomanName,omitzero"`
	BaseID              int64   `json:"baseID"`
	BaseName            string  `json:"baseName"`
	TypeDiv             int     `json:"typeDiv"`
}
```

```
{
    "baseID": 6,
    "baseName": "GBS",
    "departmentID": 5,
    "departmentName": "SATO Singapore Group 1",
    "typeDiv": 2,
    "userID": 324,
    "userName": "（TLB）TEST USER 05",
    "userNo": "555",
    "userRomanName": ""
},
```