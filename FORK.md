# Bubble Tea fork

This fork is based on Bubble Tea v2.1.1. It carries one runtime patch:
renderer shutdown waits for its goroutine to stop the ticker before a restart
resets it. A deterministic lifecycle test covers the stop/restart race.

The module path is `github.com/laith-page/bubbletea/v2`. To use a tagged fork
without a `replace`, update imports and the dependency to that path. Packages
compiled against `charm.land/bubbletea/v2` use distinct Go types and must be
updated or adapted when combined with this fork.
