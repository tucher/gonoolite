Workflow:
    1) press button on the switch n
    2) call Bind method with channel 1 + n
    3) Remember device id returned to OnBind
    4) call StartBinding with channel from (2)
    5) call Bind(0)
    6) repeat for all switches
    7) to change state of the switch call SetState(with given channel)
    8) set OnState callback to receive external updates