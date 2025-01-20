# Entity Relationship

```mermaid
erDiagram
    Style {
        string color
        string font
        string alignment
    }

    Border {
        int width
        string style
        string color
    }

    Justification {
        string horizontal
        string vertical
    }

    Screen {
        string title
        string body
        string[] questions
        Style style
    }

    Question {
        string text
        string value
        string errorMsg
    }

    Replacement {
        string color
        string[] values
    }

    Wizard {
        Screen[] screens
        bool completed
    }

    Style ||--o{ Border : references
    Style ||--o{ Justification : references
    Screen ||--o{ Style : contains
    Screen ||--o{ Replacement : interacts
    Screen ||--o{ Question : contains
    Wizard ||--o{ Screen : aggregates
    Question ||--o{ Replacement : uses
```
