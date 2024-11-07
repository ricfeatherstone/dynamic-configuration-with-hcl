# Dynamic Configuration with HCL

Quick play with Go [HCL](https://github.com/hashicorp/hcl) to understand how it could be used as a dynamic configuration 
file mechanism.

Would probably split into a hierarchy of context specific `hcl.EvalContext`s for reusability, maybe something like:
- ctx
  - stdlib
  - envfuncs
  - awsStsContext

Using fields as input to other fields e.g. `email = split(":", principal)[1]` is more complex (see second links below if
required in the future).

- [Defining Variables](https://github.com/hashicorp/hcl/blob/main/guide/go_expression_eval.rst#defining-variables)
- [Defining Functions](https://github.com/hashicorp/hcl/blob/main/guide/go_expression_eval.rst#defining-functions)
- [Interdependent Blocks](https://hcl.readthedocs.io/en/latest/go_patterns.html#interdependent-blocks)
- [Description of Multiple Phases](https://github.com/hashicorp/hcl/issues/496#issuecomment-983906130)
