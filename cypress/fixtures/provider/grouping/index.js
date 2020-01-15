Stanza(function(stanza, params) {
  const ary = [
    { x: 1, y: 1, z: 3 },
    { x: 1, y: 2, z: 4 },
    { x: 2, y: 1, z: 5 },
    { x: 2, y: 2, z: 6 },
    { x: 1, y: 2, z: 7 },
    { x: 2, y: 1, z: 8 }
  ];

  stanza.render({
    template: "stanza.html",
    parameters: {
      results: JSON.stringify(stanza.grouping(ary, ["x", "y"], "z"))
    }
  });
});
