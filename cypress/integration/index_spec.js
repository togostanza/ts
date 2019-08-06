describe('ts test', () => {
  it('Navigate to Hello Example stanza', () => {
    cy.visit('/');
    cy.contains('hello').click();
    cy.wait(50); // wait for the debunced update() function to be called

    cy.get('togostanza-hello').then($root => {
      const main = $root[0].shadowRoot.querySelector('main');
      expect(main.textContent).to.contain('Hello, world!');
    });
  })
})
