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

      expect(main.textContent).to.contain(JSON.stringify(
        [
          {
            "x_y": [
              1,
              1
            ],
            "z": [
              3
            ]
          },
          {
            "x_y": [
              1,
              2
            ],
            "z": [
              4
            ]
          },
          {
            "x_y": [
              2,
              1
            ],
            "z": [
              5
            ]
          },
          {
            "x_y": [
              2,
              2
            ],
            "z": [
              6
            ]
          },
          {
            "x_y": [
              1,
              2
            ],
            "z": [
              7
            ]
          },
          {
            "x_y": [
              2,
              1
            ],
            "z": [
              8
            ]
          }
        ]
      ));
    });
  });
});
