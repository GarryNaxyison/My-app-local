# Level Assessment Sources

The in-bot language placement tests are CEFR-style diagnostic tests. They are based on public exam formats and competency levels, but the questions are original and are not copied from proprietary exam papers.

Reference formats used:

- CEFR descriptors from the Council of Europe: https://www.coe.int/en/web/common-european-framework-reference-languages/cefr-descriptors
- English: Cambridge/IELTS-style grammar, lexis, collocation, reading and formal-register tasks.
- German: Goethe-Institut A1-C2 exam structure and practice-material task types: https://www.goethe.de/en/spr/prf.html
- Russian: TORFL/TRKI-style grammar, lexis, aspect, case government, participial/deverbal constructions, discourse markers, and formal-register precision mapped to CEFR-style placement.
- Spanish: DELE A1-C2 structure from Instituto Cervantes: https://examenes.cervantes.es/es/dele/que-es
- French: DELF/DALF public exam families from France Education International: https://www.france-education-international.fr/diplome/delf-tout-public
- Italian: CILS/CELI-style CEFR placement, using the CILS A1-C2 certification model from Universita per Stranieri di Siena: https://cils.unistrasi.it/
- Portuguese, Polish, Romanian, Ukrainian, Kazakh, Kyrgyz, Uzbek, Tajik, Tatar, Armenian, and Georgian: original CEFR-style placement tasks covering morphology, syntax, cloze completion, collocation, discourse markers, indirect speech, conditional/concessive structures, and formal-register precision. The less standardized exam families are mapped against the public CEFR descriptors rather than copied from proprietary papers.

Current production behavior:

- Every selectable learning language has a handcrafted 36-question diagnostic test.
- The old word-translation diagnostic is not used for level placement.
- Question chrome is localized to the selected bot language, while the tested material stays in the target language so the interface does not reveal the answer.
