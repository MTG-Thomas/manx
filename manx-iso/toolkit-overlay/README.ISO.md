# MANX toolkit content shipped on the ISO (root of the `/manx/` directory)
#
# Laid out so the SystemRescue live session (autorun/00-manx-setup.sh) copies
# this into /toolkit at boot:
#
#   manx/toolkit/bin/{manx,manx-tui}           the Go binaries (Tier-2 + CLI)
#   manx/toolkit/manx-iso-overlay/             Tier-1 whiptail view + helper scripts
#
# Everything in /toolkit/out/ is runtime (audit log, driver report, collected
# artifacts), produced live and never baked.
#
# The overlay scripts here are the ISO-side VIEW. Scripts must not re-implement
# verbs - the action surface lives in the Go binary (`manx <verb>` command) and
# the spec contract in bifrost-workspace docs/rescue/MANX-SPEC.md.
#
# Nothing in this directory or in the ISO may contain secrets: the SKR contract
# is enforced by the release CI's secrets scan (see SECURITY.md).
