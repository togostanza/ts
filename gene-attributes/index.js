Stanza(function(params) {
  this.query({
    endpoint: "http://togogenome.org/sparql",
    template: "stanza.rq",
    parameters: params
  }, function(rows) {
    rows.forEach(function(row) {
      row.tax_link = "http://identifiers.org/taxonomy/" + row.taxid.split(":").slice(-1)[0];
      row.refseq_link = "http://identifiers.org/refseq/" + row.refseq_label.split(":").slice(-1)[0];
    });
    this.render({
      template: "stanza.html",
      parameters: {
        gene_attributes: rows[0]
      }
    });
  });
});
