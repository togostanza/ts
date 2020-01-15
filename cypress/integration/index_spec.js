describe("hello", () => {
  it("Navigate to Hello Example stanza", () => {
    cy.visit("/");
    cy.contains("hello").click();
    cy.wait(50); // wait for the debunced update() function to be called

    cy.get("togostanza-hello").then($root => {
      const main = $root[0].shadowRoot.querySelector("main");
      expect(main.textContent).to.contain("Hello, world!");
    });
  });
});

describe("grouping", () => {
  it("processes gropuing()", () => {
    cy.visit("/");
    cy.contains("grouping").click();
    cy.wait(100); // wait for the debunced update() function to be called

    cy.get("togostanza-grouping").then($root => {
      const main = $root[0].shadowRoot.querySelector("main");
      expect(main.textContent).to.contain("Hello, grouping!");
    });
  });
});
