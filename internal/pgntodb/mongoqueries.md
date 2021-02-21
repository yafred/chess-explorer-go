### Games of a user 
```
{ $and: [ { site: "chess.com" }, { $or: [ { white: "fredo599" }, { black: "fredo599" }  ]  }  ]}
```

### Games sharing 6 same first moves
```
[{$group: {
 _id: { m01: "$m01",
  m02: "$m02",
  m03: "$m03",
  m04: "$m04",
  m05: "$m05",
  m06: "$m06",
 },
 count: { $sum: 1 }
}}, {$sort: {
  count: -1
}}, {$count: 'count'}]
```

### Time controls for a user
```
[{$match: {
  $and: [ { site: "chess.com" }, { $or: [ { white: "fredo599" }, { black: "fredo599" }  ]  }  ]
}}, {$group: {
  _id: { timecontrol: "$timecontrol"},
  count: {
    $sum: 1
  }
}}, {$sort: {
  count: -1}}]
```

### Delete games of a user (only when the opponent is not in lastgames)
```
{ $and: [ { site: "chess.com" }, 
          { $or: [ { white: "fredo599" }, { black: "fredo599" }  ]  }  ],
          { white: { $nin: [ "DBT1986", "SmoothBalance" ]} },
          { black: { $nin: [ "DBT1986", "SmoothBalance" ]} },
}
```


### Aggregate number of games for users in lastgames
```
[{$match:  {
$and: [ { site: "lichess.org" }, { $or: [ { white: "EricRosen" }, { black: "EricRosen" }  ]  }  ]
}}, {$count: 'count'}]
```

### timecontrols starting with
{ timecontrol: { $regex: /^600\+/ } }