Stanza(function(params) {
  this.query({
    endpoint: "http://togogenome.org/sparql",
    template: "stanza.rq",
    parameters: params
  }, function(rows) {
    rows.forEach(function(row) {
      row.sequence_length = row.sequence ? row.sequence.length : null;

      switch (row.fragment) {
        case 'single':
        case 'multiple':
          row.sequence_status = 'Fragment';
          break;
        default:
          row.sequence_status = 'Complete';
      }

      row.sequence_processing = row.precursor == '1' ? 'precursor' : null;
    });

    this.render({
      template: "stanza.html",
      parameters: {
        attributes: rows
      }
    });
  });
});
