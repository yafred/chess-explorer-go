### Games of a user 
```
{ $and: [ { site: "chess.com" }, { $or: [ { white: "fredo599" }, { black: "fredo599" }  ]  }  ]}
```

### Games sharing 6 same first moves
```
[{$group: {
 _id: { move01: "$move01",
  move02: "$move02",
  move03: "$move03",
  move04: "$move04",
  move05: "$move05",
  move06: "$move06",
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