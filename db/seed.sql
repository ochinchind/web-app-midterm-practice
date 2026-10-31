INSERT INTO wanted_criminals (name, nickname, reward, latitude, longitude, last_known_note, difficulty, armed, has_crew, escape_master) VALUES
 ('Jonas Varek',  'Red Comet',      25000,  34.0522, -118.2437, 'Seen at Luna Port docking bay 4', 7, TRUE,  TRUE,  TRUE),
 (NULL,           'The Ghost',      80000,  NULL,    NULL,      'Last seen near the Kuiper Belt relay station', 10, TRUE, FALSE, TRUE),
 ('Mara Quill',   'Void Widow',     12000, -33.8688,  151.2093, NULL, 5, FALSE, TRUE,  FALSE),
 ('Tobias Krane', 'Iron Jaw',        5000,  51.5074,   -0.1278, 'Bar fight on Titan Station', 3, TRUE,  FALSE, FALSE),
 (NULL,           'Dust Devil',      7500,  NULL,    NULL,      'Rumored to hide on Mars colony Olympus', 4, FALSE, FALSE, TRUE),
 ('Lyra Venn',    'Nova Queen',     45000,  35.6762,  139.6503, 'Boarded a freighter bound for Ceres', 8, TRUE,  TRUE,  FALSE),
 ('Ruk Dalen',    'Scrap King',      3000,  NULL,    NULL,      'Salvage yard on Europa', 2, FALSE, FALSE, FALSE);

INSERT INTO hunt_logs (criminal_id, success, notes, created_at) VALUES
 (1, FALSE, 'Target escaped through an asteroid field.',           now() - interval '10 days'),
 (1, FALSE, 'Lost the trail after a hyperspace jump.',             now() - interval '8 days'),
 (1, TRUE,  'Cornered at Luna Port, but he later bribed his way out.', now() - interval '5 days'),
 (2, FALSE, 'Only found an empty ship and a decoy transponder.',   now() - interval '9 days'),
 (2, FALSE, 'Ambushed by unknown gunmen near the relay station.',  now() - interval '4 days'),
 (3, TRUE,  'Captured without resistance; crew surrendered.',      now() - interval '7 days'),
 (4, TRUE,  'Arrested at the Titan Station bar.',                  now() - interval '6 days'),
 (4, FALSE, 'False tip, target was never there.',                  now() - interval '3 days'),
 (5, FALSE, 'Sandstorm on Mars covered his tracks.',               now() - interval '2 days'),
 (6, FALSE, 'Freighter crew refused to cooperate.',                now() - interval '36 hours'),
 (6, FALSE, 'Gunfight in the cargo bay, target slipped away.',     now() - interval '20 hours'),
 (7, TRUE,  'Captured at the Europa salvage yard.',                now() - interval '1 hour');
